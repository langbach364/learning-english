#!/bin/bash

FILE_TO_WATCH="./trans.txt"
TRANSLATED_FILE="./trans_ed.txt"
LAST_TRANSLATED_FILE="./last_translated.txt"
TRANSLATION_MODEL="auto" 

if [ ! -f "$FILE_TO_WATCH" ]; then
    touch "$FILE_TO_WATCH"
fi

translate_file() {
    local input_file=$1
    local output_file=$2

    trans -b -i "$input_file" -o "$output_file" -s en -t vi -e "$TRANSLATION_MODEL"
    cp "$input_file" "$LAST_TRANSLATED_FILE"
}

inotifywait -m -e close_write --format '%w%f' "$FILE_TO_WATCH" | while read -r NEW_FILE; do
    if [ "$NEW_FILE" = "$FILE_TO_WATCH" ]; then
        if ! cmp -s "$FILE_TO_WATCH" "$LAST_TRANSLATED_FILE"; then
            echo "Phát hiện nội dung mới. Bắt đầu dịch file..."
            translate_file "$NEW_FILE" "$TRANSLATED_FILE"
            echo "Dịch xong. Kết quả được lưu trong $TRANSLATED_FILE"
        else
            echo "Không có thay đổi, bỏ qua dịch thuật."
        fi
    fi
done
