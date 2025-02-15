package abhash

import "crypto/sha256"

// 🔥 ABHash implements a hash function that:
//   - Takes:
//     – data: input byte slice ([]byte)
//     – sohp: expected size of each part for hashing (default is 2)
//     – hpa: number of hash parts (default is 2)
//   - Splits the input data into hpa parts so that the first hpa-1 parts have a length of sohp,
//     and the last part gets the remaining bytes (it may be larger than sohp).
//   - Computes a fixed-size token for each part using XOR-folding and bitwise inversion.
//     This ensures that changes in data only affect the token of the part where the change occurred.
func ABHash(data []byte, sohp int, hpa int) []byte {
	// 💡 Set default values if invalid parameters are provided
	if sohp <= 0 {
		sohp = 2
	}
	if hpa <= 0 {
		hpa = 2
	}

	// 📋 Split data into hpa parts.
	// The first hpa-1 parts have exactly sohp bytes (or fewer if data is short, using remaining bytes).
	// The last part gets all the remaining data.
	parts := make([][]byte, hpa)
	for i := 0; i < hpa-1; i++ {
		start := i * sohp
		end := start + sohp

		if start > len(data) {
			start = len(data)
		}

		if end > len(data) { // if data is smaller than required for the current part
			end = len(data)
		}
		parts[i] = data[start:end]
	}
	// The last part gets all remaining data
	if (hpa-1)*sohp < len(data) {
		parts[hpa-1] = data[(hpa-1)*sohp:]
	} else {
		// Add padding to the last part if it's smaller than sohp bytes
		paddingSize := sohp - len(parts[hpa-1])
		padding := make([]byte, paddingSize)
		parts[hpa-1] = append(parts[hpa-1], padding...)
	}

	// 🤗 Generate tokens for each part
	tokens := make([][]byte, hpa)
	for i, part := range parts {
		tokens[i] = generateToken(part, sohp)
	}

	// 🔄 Assemble the final hash by concatenating tokens in the same order as data parts
	result := make([]byte, 0, hpa*sohp)
	for _, token := range tokens {
		result = append(result, token...)
	}

	return result
}

// SECRET is a secret key used in generateToken for hashing with salt
var SECRET []byte

// GgenerateToken generates a token of a fixed size from a given data part using XOR-folding and bitwise inversion.
//
// The function takes two parameters:
// - part: a slice of bytes representing the data part from which the token will be generated.
// - size: an integer representing the desired size of the token.
//
// The function returns a slice of bytes representing the generated token.
//
// The function works as follows:
// 1. Creates a new SHA256 hash object.
// 2. Writes the data part and the SECRET (if not nil) to the hash object.
// 3. Retrieves the hash sum from the hash object.
// 4. Initializes a token slice of the specified size.
// 5. Performs XOR-folding by adding each byte of the hashed data to the corresponding index in the token slice, modulo the token size.
// 6. Applies bitwise inversion to each byte in the token slice.
// 7. Returns the generated token.
func generateToken(part []byte, size int) []byte {
	// Create a hash object
	hash := sha256.New()

	// Add the data part to the hash
	hash.Write(part)

	if SECRET == nil {
		SECRET = []byte{} // Assign an empty array if SECRET was not set before
	}

	hash.Write(SECRET)

	// Get the hash sum
	hashedData := hash.Sum(nil)
	token := make([]byte, size)

	// ⏱️ XOR-folding: add bytes of the hash modulo the token size
	for i, b := range hashedData {
		token[i%size] ^= b
	}
	// Apply bitwise inversion to each byte of the token
	for i := 0; i < size; i++ {
		token[i] = ^token[i]
	}

	return token
}
