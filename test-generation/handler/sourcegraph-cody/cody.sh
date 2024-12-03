#!/bin/bash

cd "$HOME"/Documents/learning-english/test-generation/handler || exit

DATA="./sourcegraph-cody/data.txt"
PROMPT=$(cat ./sourcegraph-cody/prompt.txt)
MODEL=$(cat ./sourcegraph-cody/model.txt)
FORM="./sourcegraph-cody/form.txt"
FORM_OUTPUT="./sourcegraph-cody/form-output.txt"

ANSWER_DIR="./sourcegraph-cody"
mkdir -p "$ANSWER_DIR"

ANSWER_FILE="$ANSWER_DIR/answer.txt"

cody chat --model "$MODEL" --context-file ./$DATA ./$FORM ./$FORM_OUTPUT --message "$PROMPT"  >"$ANSWER_FILE"

echo "Câu trả lời đã được lưu vào: $ANSWER_FILE"
exit 0
