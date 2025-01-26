# read-only

```
git pull
apt install git-annex / brew install git-annex
./get_testfiles.sh
```

# read/write access

```
git pull
apt install git-annex / brew install git-annex
WEBDAV_USERNAME=... WEBDAV_PASSWORD=... git annex enableremote cave.servium.ch
git annex sync
git annex get *
```

# add a file

```
git annex add *
git commit
git annex copy * --to cave.servium.ch
```
