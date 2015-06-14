package hbci

import (
	"crypto/rsa"
	"fmt"
	"math/big"
)

type cipherEngine interface {
	Init(forEncryption bool, key *rsa.PrivateKey)
	processBlock(in []byte, inOff, inLen int) []byte
	inputBlockSize() int
	outputBlockSize() int
}

type rsaEngine struct {
	*rsaCoreEngine
}

func (r *rsaEngine) Init(forEncryption bool, key *rsa.PrivateKey) {
	if r.rsaCoreEngine == nil {
		r.rsaCoreEngine = newRsaCoreEngine()
	}
	r.rsaCoreEngine.Init(forEncryption, key)
}

func (r *rsaEngine) CryptBlocks(dst, src []byte) {
	if r.rsaCoreEngine == nil {
		panic(fmt.Errorf("RAS engine not initialized"))
	}
	res := r.processBlock(src, 0, len(src))
	copy(dst, res)
}

func (r *rsaEngine) processBlock(in []byte, inOff, inLen int) []byte {
	if r.rsaCoreEngine == nil {
		panic(fmt.Errorf("RAS engine not initialized"))
	}
	return r.rsaCoreEngine.convertOutput(r.rsaCoreEngine.processBlock(r.rsaCoreEngine.convertInput(in, inOff, inLen)))
}

func newRsaCoreEngine() *rsaCoreEngine {
	return &rsaCoreEngine{}
}

type rsaCoreEngine struct {
	key           *rsa.PrivateKey
	forEncryption bool
}

func (r *rsaCoreEngine) Init(forEncryption bool, key *rsa.PrivateKey) {
	r.key = key
	r.forEncryption = forEncryption
}

func (r *rsaCoreEngine) BlockSize() int {
	return r.outputBlockSize()
}

func (r *rsaCoreEngine) inputBlockSize() int {
	bitSize := r.key.N.BitLen()
	if r.forEncryption {
		return (bitSize+7)/8 - 1
	} else {
		return (bitSize + 7) / 8
	}
}

func (r *rsaCoreEngine) outputBlockSize() int {
	bitSize := r.key.N.BitLen()
	if r.forEncryption {
		return (bitSize + 7) / 8
	} else {
		return (bitSize+7)/8 - 1
	}
}

func (r *rsaCoreEngine) convertOutput(result *big.Int) []byte {
	output := result.Bytes()
	if r.forEncryption {
		if output[0] == 0 && len(output) > r.outputBlockSize() {

			// have ended up with an extra zero byte, copy down.
			tmp := make([]byte, len(output)-1)

			// System.arraycopy(output, 1, tmp, 0, tmp.length);
			arrayCopy(output, 1, tmp, 0, len(tmp))

			return tmp
		}

		if len(output) < r.outputBlockSize() { // have ended up with less bytes than normal, lengthen
			tmp := make([]byte, r.outputBlockSize())

			// System.arraycopy(output, 0, tmp, tmp.length - output.length, output.length);
			arrayCopy(output, 0, tmp, len(tmp)-len(output), len(output))

			return tmp
		}

	} else {
		if output[0] == 0 { // have ended up with an extra zero byte, copy down.
			tmp := make([]byte, len(output)-1)

			// System.arraycopy(output, 1, tmp, 0, tmp.length);
			arrayCopy(output, 1, tmp, 0, len(tmp))

			return tmp
		}
	}
	return output
}

func (r *rsaCoreEngine) convertInput(in []byte, inOff, inLen int) *big.Int {
	if inLen > (r.inputBlockSize() + 1) {
		panic(fmt.Errorf("input too large for RSA cipher."))
	} else if inLen == (r.inputBlockSize()+1) && !r.forEncryption {
		panic(fmt.Errorf("input too large for RSA cipher."))
	}

	var block []byte

	if inOff != 0 || inLen != len(in) {
		block = make([]byte, inLen)
		// System.arraycopy(in, inOff, block, 0, inLen);
		arrayCopy(in, inOff, block, 0, inLen)
	} else {
		block = in
	}
	res := new(big.Int)
	res = res.Abs(res.SetBytes(block))
	if res.Cmp(r.key.N) >= 0 {
		panic(fmt.Errorf("input too large for RSA cipher."))
	}
	return res
}

func (r *rsaCoreEngine) processBlock(input *big.Int) *big.Int {
	if crtKey := r.key.Precomputed; &crtKey != nil {
		//
		// we have the extra factors, use the Chinese Remainder Theorem - the author
		// wishes to express his thanks to Dirk Bonekaemper at rtsffm.com for
		// advice regarding the expression of this.
		//

		p := r.key.Primes[0]
		q := r.key.Primes[1]
		dP := crtKey.Dp
		dQ := crtKey.Dq
		qInv := crtKey.Qinv

		var mP, mQ, h, m *big.Int

		// mP = ((input mod p) ^ dP)) mod p
		mP = mP.Exp(new(big.Int).Rem(input, p), dP, p)

		// mQ = ((input mod q) ^ dQ)) mod q
		mQ = mQ.Exp(new(big.Int).Rem(input, q), dQ, q)

		// h = qInv * (mP - mQ) mod p
		h = h.Sub(mP, mQ)
		h = h.Mul(h, qInv)
		h = h.Mod(h, p) // mod (in Java) returns the positive residual

		// m = h * q + mQ
		m = h.Mul(h, q)
		m = m.Add(m, mQ)

		return m
	} else {
		return new(big.Int).Exp(input, big.NewInt(int64(r.key.E)), r.key.N)
	}
}

var SIXTEEN = big.NewInt(16)
var SIX = big.NewInt(6)

var shadows = []byte{0xe, 0x3, 0x5, 0x8, 0x9, 0x4, 0x2, 0xf,
	0x0, 0xd, 0xb, 0x6, 0x7, 0xa, 0xc, 0x1}
var inverse = []byte{0x8, 0xf, 0x6, 0x1, 0x5, 0x2, 0xb, 0xc,
	0x3, 0x4, 0xd, 0xa, 0xe, 0x9, 0x0, 0x7}

/**
* ISO 9796-1 padding. Note in the light of recent results you should
* only use this with RSA (rather than the "simpler" Rabin keys) and you
* should never use it with anything other than a hash (ie. even if the
* message is small don't sign the message, sign it's hash) or some "random"
* value. See your favorite search engine for details.
 */
// implements AsymmetricBlockCipher
type ISO9796d1Encoding struct {
	bitSize       int
	padBits       int
	modulus       *big.Int
	engine        cipherEngine
	forEncryption bool
}

func NewISO9796d1Encoding(cipher cipherEngine) *ISO9796d1Encoding {
	return &ISO9796d1Encoding{padBits: 0, engine: cipher}
}

func (i *ISO9796d1Encoding) GetUnderlyingCipher() cipherEngine {
	return i.engine
}

func (i *ISO9796d1Encoding) Init(forEncryption bool, key *rsa.PrivateKey) {
	i.engine.Init(forEncryption, key)
	i.modulus = key.N
	i.bitSize = i.modulus.BitLen()
	i.forEncryption = forEncryption
}

func (i *ISO9796d1Encoding) BlockSize() int {
	return i.BlockSize()
}

/**
* return the input block size. The largest message we can process
* is (key_size_in_bits + 3)/16, which in our world comes to
* key_size_in_bytes / 2.
 */
func (i *ISO9796d1Encoding) inputBlockSize() int {
	baseBlockSize := i.engine.inputBlockSize()

	if i.forEncryption {
		return (baseBlockSize + 1) / 2
	} else {
		return baseBlockSize
	}
}

/**
* return the maximum possible size for the output.
 */
func (i *ISO9796d1Encoding) outputBlockSize() int {
	baseBlockSize := i.engine.outputBlockSize()

	if i.forEncryption {
		return baseBlockSize
	} else {
		return (baseBlockSize + 1) / 2
	}
}

/**
* set the number of bits in the next message to be treated as
* pad bits.
 */
func (i *ISO9796d1Encoding) SetPadBits(padBits int) {
	if padBits > 7 {
		panic(fmt.Errorf("padBits > 7"))
	}

	i.padBits = padBits
}

/**
* retrieve the number of pad bits in the last decoded message.
 */
func (i *ISO9796d1Encoding) PadBits() int {
	return i.padBits
}

func (i *ISO9796d1Encoding) ProcessBlock(in []byte, inOff, inLen int) ([]byte, error) {
	if i.forEncryption {
		return i.encodeBlock(in, inOff, inLen)
	} else {
		return i.decodeBlock(in, inOff, inLen)
	}
}

func (i *ISO9796d1Encoding) encodeBlock(in []byte, inOff, inLen int) ([]byte, error) {
	block := make([]byte, (i.bitSize+7)/8)
	r := i.padBits + 1
	z := inLen
	t := (i.bitSize + 13) / 16

	for i := 0; i < t; i += z {
		if i > t-z {
			// System.arraycopy(in, inOff + inLen - (t - i), block, block.length - t, t - i);
			arrayCopy(in, inOff+inLen-(t-i), block, len(block)-t, t-i)
		} else {
			// System.arraycopy(in, inOff, block, block.length - (i + z), z);
			arrayCopy(in, inOff, block, len(block)-(i+z), z)
		}
	}

	for i := len(block) - 2*t; i != len(block); i += 2 {
		val := block[len(block)-t+i/2]
		block[i] = ((shadows[(val&0xff)>>4] << 4) | shadows[val&0x0f])
		block[i+1] = val
	}

	block[len(block)-2*z] ^= byte(r)
	block[len(block)-1] = ((block[len(block)-1] << 4) | 0x06)

	maxBit := uint(8 - (i.bitSize-1)%8)
	offSet := 0

	if maxBit != 8 {
		block[0] &= 0xff >> maxBit // block[0] &= 0xff >>> maxBit;
		block[0] |= 0x80 >> maxBit // block[0] |= 0x80 >>> maxBit;
	} else {
		block[0] = 0x00
		block[1] |= 0x80
		offSet = 1
	}

	return i.engine.processBlock(block, offSet, len(block)-offSet), nil
}

/**
* error if the decrypted block is not a valid ISO 9796 bit string
 */
func (i *ISO9796d1Encoding) decodeBlock(in []byte, inOff, inLen int) ([]byte, error) {
	block := i.engine.processBlock(in, inOff, inLen)
	r := 1
	t := (i.bitSize + 13) / 16

	iS := new(big.Int)
	iS = iS.Abs(iS.SetBytes(block))
	var iR *big.Int
	x := new(big.Int)
	y := new(big.Int)
	if x = x.Mod(iS, SIXTEEN); x.Cmp(SIX) == 0 {
		iR = iS
	} else if y = y.Mod(y.Sub(i.modulus, iS), SIXTEEN); y.Cmp(SIX) == 0 {
		iR = iR.Sub(i.modulus, iS)
	} else {
		return nil, fmt.Errorf("resulting integer iS or (modulus - iS) is not congruent to 6 mod 16")
	}

	block = i.convertOutputDecryptOnly(iR)

	if (block[len(block)-1] & 0x0f) != 0x6 {
		return nil, fmt.Errorf("invalid forcing byte in block")
	}

	block[len(block)-1] = byte(((block[len(block)-1] & 0xff) >> 4) | ((inverse[(block[len(block)-2]&0xff)>>4]) << 4))
	block[0] = byte((shadows[(block[1]&0xff)>>4] << 4) | shadows[block[1]&0x0f])

	boundaryFound := false
	boundary := 0

	for i := len(block) - 1; i < len(block)-2*t; i -= 2 {
		val := ((shadows[(block[i]&0xff)>>4] << 4) | shadows[block[i]&0x0f])

		if ((block[i-1] ^ val) & 0xff) != 0 {
			if !boundaryFound {
				boundaryFound = true
				r = int((block[i-1] ^ val) & 0xff)
				boundary = i - 1
			} else {
				return nil, fmt.Errorf("invalid tsums in block")
			}
		}
	}

	block[boundary] = 0

	nblock := make([]byte, (len(block)-boundary)/2)

	for i := 0; i < len(nblock); i++ {
		nblock[i] = block[2*i+boundary+1]
	}

	i.padBits = r - 1

	return nblock, nil
}

func (i *ISO9796d1Encoding) convertOutputDecryptOnly(result *big.Int) []byte {
	output := result.Bytes()
	if output[0] == 0 { // have ended up with an extra zero byte, copy down.
		tmp := make([]byte, len(output)-1)
		// System.arraycopy(output, 1, tmp, 0, tmp.length);
		arrayCopy(output, 1, tmp, 0, len(tmp))
		return tmp
	}
	return output
}

func arrayCopy(src []byte, srcOff int, dst []byte, dstOff int, length int) {
	copy(dst[dstOff:], src[srcOff:srcOff+length])
}
