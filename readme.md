# dbdump

Retrieve forgotten password from DBeaver.

### Build instructions

```bash
gh repo clone gregzakh/dbdump
cd dbdump
go build
```

### Usage examples

```bash
# compatible with default installation
./dbdump
# useful with custom paths such as snap
./dbdump -c ~/snap/dbeaver-ce/418/.local/share
```
