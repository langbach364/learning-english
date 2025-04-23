#!/bin/bash

# cd ./learning-english/test-generation/handler/ || exit
# cd ./sourcegraph-cody/test-generator/ || exit

DATA="./data.txt"
PROMPT=$(cat ./prompt.txt)
MODEL=$(cat ./model.txt)

FORM="./form.txt"
FORM_OUTPUT="./form-output.txt"

ANSWER_DIR="./"
mkdir -p "$ANSWER_DIR"

ANSWER_FILE="$ANSWER_DIR/answer.txt"

cody chat --model "$MODEL" --context-file ./$DATA ./$FORM ./$FORM_OUTPUT --message "$PROMPT"  >"$ANSWER_FILE"

echo "Câu trả lời đã được lưu vào: $ANSWER_FILE"
exit 0
