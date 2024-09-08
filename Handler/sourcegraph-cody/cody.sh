#!/bin/bash

DATA=$(cat ./sourcegraph-cody/data.txt)
PROMPT=$(cat ./sourcegraph-cody/prompt.txt)
MODEL=$(cat ./sourcegraph-cody/model.txt)

COMBINED_PROMPT="${DATA}
${PROMPT}"


ANSWER_DIR="./sourcegraph-cody"
mkdir -p "$ANSWER_DIR"

ANSWER_FILE="$ANSWER_DIR/answer.txt"

cody chat --model "$MODEL" -m "$COMBINED_PROMPT" > "$ANSWER_FILE"

echo "Câu trả lời đã được lưu vào: $ANSWER_FILE"
