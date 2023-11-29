# hashifuse-go

*draft: v0.0.1a*

## warning
  - filesystem is filled once on startup and persist in memory with all secrets content
  - not secure
  - not production ready

## overview
hashifuse-go is a single binary to read-only mount hashicorp vault's engine as local filesystem via fuse

## build
    go build .

## short help
	$ ./hashifuse-go --help
	Usage of ./hashifuse-go:
	  -debug_fuse
	    	enable debug for fuse operations
	  -mnt_dir string
	    	set mount dir (default "/mnt/vaultfs")
	  -secret_engine_path string
	    	set the name of secret engine to mount (default "kv")
	  -vault_addr string
	    	set hashicorp vault address with vaultfs scheme format (default "vault+http://127.0.0.1:8200")

## usage
    VAULT_TOKEN=<vault_token> ./hashifuse-go -vault_addr "<vaultfs_schema>://<vault-addr>" -secret_engine_path "<secret_engine>" -mnt_dir "<mount_directory>"
where:
  - `vault_token` is your hashicorp vault token
  - `vaultfs_schema` is one of [vaultfs schemas](https://github.com/hairyhenderson/go-fsimpl/tree/main#supported-filesystems)
  - `secret_engine` is secret engine `Path` without trailing slash (ex., `secret`)
  - `mount_directory` is the mount directory, must exist
## example
### mount
    VAULT_TOKEN='hvs.9Wzzyqq4T1AnPME3B7g3Vero' ./hashifuse-go -vault_addr "vault+https://vault.example.com" -secret_engine_path "secret" -mnt_dir "/mnt/vaultfs"
### list secrets
    ls -lR /mnt/vaultfs

## unmount
    fusermount -u <mount_directory>

## docs
  - https://github.com/hairyhenderson/go-fsimpl/tree/main/vaultfs
  - https://pkg.go.dev/github.com/hanwen/go-fuse/v2/fs#section-readme
  - https://pkg.go.dev/io/fs
  - code examples on the web

## inspired by
  - https://www.hashicorp.com/resources/hashicorp-apis-via-hashifuse

## TODO
  - remove trailing slash in `secret_engine` argument, if exists
  - make it read-write
  - cleanup code, workaround exceptions and so on
