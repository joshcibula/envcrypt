# envcrypt

> Lightweight `.env` file encryption and sharing tool with [age](https://github.com/FiloSottile/age) encryption backend.

---

## Installation

```bash
go install github.com/yourusername/envcrypt@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envcrypt.git
cd envcrypt && go build -o envcrypt .
```

---

## Usage

**Encrypt a `.env` file:**

```bash
envcrypt encrypt -i .env -o .env.age -r age1ql3z7hjy...
```

**Decrypt a `.env` file:**

```bash
envcrypt decrypt -i .env.age -o .env -k ~/.age/key.txt
```

**Generate a new age key pair:**

```bash
envcrypt keygen
```

Share the encrypted `.env.age` file safely in version control or over untrusted channels. Only recipients with the corresponding private key can decrypt it.

---

## How It Works

`envcrypt` wraps the [age](https://age-encryption.org) encryption standard to provide a simple CLI for encrypting and decrypting `.env` files. Keys can be recipient-based (public key) or passphrase-based.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)