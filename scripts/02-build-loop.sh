#/usr/bin/env bash
CMD="go build -x"
which colorgo >/dev/null
if [ $? -eq 0 ]; then
    CMD="colorgo build -x"
fi

while true; do
    gate || break
    clear
    $CMD
done

