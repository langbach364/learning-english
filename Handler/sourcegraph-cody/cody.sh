#!/bin/bash

DATA=$(cat learning-english/Handler/sourcegraph-cody/data.txt)
PROMPT=$(cat learning-english/Handler/sourcegraph-cody/prompt.txt)
MODEL=$(cat learning-english/Handler/sourcegraph-cody/model.txt)

COMBINED_PROMPT="${DATA}
${PROMPT}"


ANSWER_DIR="learning-english/Handler/sourcegraph-cody"
mkdir -p "$ANSWER_DIR"

ANSWER_FILE="$ANSWER_DIR/answer.txt"

cody chat --model "$MODEL" -m "$COMBINED_PROMPT" > "$ANSWER_FILE"

echo "Câu trả lời đã được lưu vào: $ANSWER_FILE"
