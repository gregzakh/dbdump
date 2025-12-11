```bash
#!/usr/bin/env bash

find ~/.local -type f -regextype posix-extended -regex '.*/(cred.*|data.*)\.json' -print0 |
xargs -0 -n1 sh -c '
for f do
   case "$f" in
      *cred*.json) openssl aes-128-cbc -d -K babb4a9f774ab853c96c2d653dfe544a \
                   -iv 00000000000000000000000000000000 -in "$f" | cut -c17- |
                   yq ".[][\"#connection\"] | \"\(.user):\(.password)\"" ;;
      *data*.json) yq ".connections[].configuration.url" < "$f" ;;
   esac
done
' sh | awk '
/^jdbc/{jdbc[++j]=$0} !/^jdbc/{auth[++a]=$0}
END{for(i=1;i<=j;i++){split(auth[i],b,":");print jdbc[i] " | " b[1] " | " b[2]}}
' | column -t
```
