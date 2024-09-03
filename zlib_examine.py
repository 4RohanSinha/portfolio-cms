import zlib
import sys

def dump_zlib_file(input_file):
    try:
        # Read the zlib compressed file
        with open(input_file, 'rb') as f:
            compressed_data = f.read()
        
        # Decompress the data
        decompressed_data = zlib.decompress(compressed_data)
        
        # Output the decompressed data to stdout
        sys.stdout.buffer.write(decompressed_data)
    
    except FileNotFoundError:
        print(f"File '{input_file}' not found.", file=sys.stderr)
    except zlib.error as e:
        print(f"Error decompressing file: {e}", file=sys.stderr)
    except Exception as e:
        print(f"An unexpected error occurred: {e}", file=sys.stderr)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python dump_zlib.py <input_file>", file=sys.stderr)
    else:
        input_file = sys.argv[1]
        dump_zlib_file(input_file)

