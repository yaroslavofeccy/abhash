# ABHash

**ABHash** is a lightweight, high-performance hash function library written in Go. It is designed to split your input data into multiple segments and compute a unique fixed-size token for each segment using XOR folding and bitwise inversion. This approach allows you to identify issues in specific segments by ensuring that changes in the input only affect the corresponding token.

## Features

- **Modular Hashing:** Independently hashes segments of the input data.
- **Configurable Parameters:** Set the size of each hash part and the number of parts.
- **Efficient Performance:** Utilizes XOR folding and bitwise inversion for fast computations.
- **Incremental Verification:** Only modified segments produce different tokens, making it easy to isolate changes.

## Installation

To use **abhash** in your project, initialize your module and add the dependency:

```bash
go get github.com/yaroslavofeccy/abhash
```
Usage
Import the module in your Go code:

```go
import "github.com/yourusername/abhash"
```

Here’s a quick example to demonstrate its usage:
```go
package main

import (
	"fmt"
	"github.com/yourusername/abhash"
)

func main() {
	// Input data to be hashed
	data := []byte("This is an example of data to hash using abhash!")
	
	// Define parameters
	sohp := 4 // Expected size of each token part (in bytes)
	hpa := 3  // Number of hash parts

	// Compute the hash
	hash := abhash.FastHash(data, sohp, hpa)

	// Output the results
	fmt.Printf("Input Data: %q\n", data)
	fmt.Printf("Hash: %x\n", hash)
}
```

How It Works
Data Splitting:
The input data is divided into hpa segments. The first hpa - 1 segments have a fixed size (sohp bytes), while the last segment contains the remaining bytes.

Token Generation:
For each segment, the library computes a token of size sohp using a two-step process:

XOR Folding: Each byte of the segment contributes to the token via XOR folding.
Bitwise Inversion: Each resulting byte in the token is then inverted using the bitwise NOT operator (^).
Hash Assembly:
The final hash is constructed by concatenating the tokens of all segments in their original order

## Testing
To run tests for abhash, use the following command from the module's root directory:

```bash
go test ./...
```

The tests verify that:
- The hash function is consistent for identical inputs.
- Changes in one part of the data affect only the corresponding token.
- The function properly handles input data shorter than the expected length.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
