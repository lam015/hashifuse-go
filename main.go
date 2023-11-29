package main

import (
	"context"
	"log"
    "os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
    
    "flag"
)

// flags vars
var (
    vault_addr string
    secret_engine_path string
    mnt_dir string
    debug_fuse bool
)

// map for files and its content
var vault_secrets = map[string]string {}

// inMemoryFS is the root of the tree
type inMemoryFS struct {
	fs.Inode
}

// Ensure that we implement NodeOnAdder
var _ = (fs.NodeOnAdder)((*inMemoryFS)(nil))

// OnAdd is called on mounting the file system. Use it to populate
// the file system tree.
func (root *inMemoryFS) OnAdd(ctx context.Context) {
	for name, content := range vault_secrets {
		dir, base := filepath.Split(name)

		p := &root.Inode

		// Add directories leading up to the file.
		for _, component := range strings.Split(dir, "/") {
			if len(component) == 0 {
				continue
			}
			ch := p.GetChild(component)
			if ch == nil {
				// Create a directory
				ch = p.NewPersistentInode(ctx, &fs.Inode{},
					fs.StableAttr{Mode: syscall.S_IFDIR})
				// Add it
				p.AddChild(component, ch, true)
			}

			p = ch
		}

		// Make a file out of the content bytes. This type
		// provides the open/read/flush methods.
		embedder := &fs.MemRegularFile{
			Data: []byte(content),
		}

		// Create the file. The Inode must be persistent,
		// because its life time is not under control of the
		// kernel.
		child := p.NewPersistentInode(ctx, embedder, fs.StableAttr{})

		// And add it
		p.AddChild(base, child, true)
	}
}

// This demonstrates how to build a file system in memory. The
// read/write logic for the file is provided by the MemRegularFile type.
func main() {
    
    flag.StringVar(&vault_addr, "vault_addr", "vault+http://127.0.0.1:8200", "set hashicorp vault address with vaultfs scheme format")
    flag.StringVar(&secret_engine_path, "secret_engine_path", "kv", "set the name of secret engine to mount")
    flag.StringVar(&mnt_dir, "mnt_dir", "/mnt/vaultfs", "set mount dir")
    flag.BoolVar(&debug_fuse, "debug_fuse", false, "enable debug for fuse operations")
    flag.Parse()

    mnt_dir_stat, err := os.Stat(mnt_dir)
    if os.IsNotExist(err) {
        log.Fatal("mount directory \"" + mnt_dir + "\" does not exist.")
    }
    log.Println(mnt_dir_stat)
    
    // create fuse fs in memory with content filled from hashicorp vault
    create_vaultfs(vault_addr, secret_engine_path)

	root := &inMemoryFS{}
	server, err := fs.Mount(mnt_dir, root, &fs.Options{
		MountOptions: fuse.MountOptions{Debug: debug_fuse},
	})
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Mounted on %s", mnt_dir)
	log.Printf("Unmount by calling 'fusermount -u %s'", mnt_dir)

	// Wait until unmount before exiting
	server.Wait()
}
