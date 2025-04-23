#!/bin/bash

# cd ./learning-english/chat-cody/handler/ || exit
cd ./sourcegraph-cody/chat-cody/ || exit

LOG_FILE="./cody_run.log"
echo "==================================================" >>"$LOG_FILE"
echo "Script started at: $(date)" >>"$LOG_FILE"

QUESTION_DIR="./question"
ANSWER_DIR="./answer"
mkdir -p "$QUESTION_DIR"
mkdir -p "$ANSWER_DIR"
echo "Checked/Created directories: $QUESTION_DIR, $ANSWER_DIR" >>"$LOG_FILE"

MODEL_FILE="./model.txt"
PROMPT_FILE="./prompt.txt"
DATA_FILE="./data.txt"

if [ ! -f "$MODEL_FILE" ]; then
  echo "ERROR: Model file not found: $MODEL_FILE" | tee -a "$LOG_FILE"
  exit 1
fi
if [ ! -f "$PROMPT_FILE" ]; then
  echo "ERROR: Prompt file not found: $PROMPT_FILE" | tee -a "$LOG_FILE"
  exit 1
fi

MODEL=$(cat "$MODEL_FILE")
PROMPT=$(cat "$PROMPT_FILE")

echo "Model read from $MODEL_FILE: $MODEL" >>"$LOG_FILE"
echo "Prompt read from $PROMPT_FILE" >>"$LOG_FILE"

ANSWER_FILE_BASE="${ANSWER_DIR}/answer.txt"
ANSWER_FILE="$ANSWER_FILE_BASE"
if [ -f "$ANSWER_FILE_BASE" ]; then
  COUNT=0
  for f in "${ANSWER_DIR}"/answer\(*\).txt; do
    if [[ -e "$f" ]]; then
      num=$(echo "$f" | sed -n 's/.*answer(\([0-9]*\)).txt/\1/p')
      if [[ "$num" =~ ^[0-9]+$ && "$num" -gt "$COUNT" ]]; then
        COUNT=$num
      fi
    fi
  done
  NEXT_COUNT=$((COUNT + 1))
  ANSWER_FILE="${ANSWER_DIR}/answer(${NEXT_COUNT}).txt"
  echo "Base answer file $ANSWER_FILE_BASE exists. Found max count: $COUNT. Next answer file will be: $ANSWER_FILE" >>"$LOG_FILE"
else
  echo "Base answer file $ANSWER_FILE_BASE does not exist. Using it as the answer file." >>"$LOG_FILE"
fi

declare -a CONTEXT_ARGS

if [ -f "$DATA_FILE" ]; then
  CONTEXT_ARGS+=("--context-file" "$DATA_FILE")
  echo "Adding data file to context args: $DATA_FILE" >>"$LOG_FILE"
fi

echo "DEBUG: Listing all files in $QUESTION_DIR:" >>"$LOG_FILE"
ls -la "$QUESTION_DIR" >>"$LOG_FILE"

echo "Finding latest question file for context:" >>"$LOG_FILE"
LATEST_QUESTION="${QUESTION_DIR}/question.txt"
LATEST_QUESTION_NUM=0

if [ -f "$LATEST_QUESTION" ]; then
  echo "  Found base question file: $LATEST_QUESTION" >>"$LOG_FILE"
else
  LATEST_QUESTION=""
fi

for qfile in "${QUESTION_DIR}"/question\(*\).txt; do
  if [[ -e "$qfile" ]]; then
    num=$(echo "$qfile" | sed -n 's/.*question(\([0-9]*\)).txt/\1/p')
    if [[ "$num" =~ ^[0-9]+$ && "$num" -gt "$LATEST_QUESTION_NUM" ]]; then
      LATEST_QUESTION_NUM=$num
      LATEST_QUESTION="$qfile"
    fi
  fi
done

if [ -n "$LATEST_QUESTION" ] && [ -f "$LATEST_QUESTION" ]; then
  CONTEXT_ARGS+=("--context-file" "$LATEST_QUESTION")
  echo "  Added latest question file to context: $LATEST_QUESTION" >>"$LOG_FILE"
else
  echo "  WARNING: No question file found to add to context!" >>"$LOG_FILE"
fi

echo "DEBUG: Listing all files in $ANSWER_DIR:" >>"$LOG_FILE"
ls -la "$ANSWER_DIR" >>"$LOG_FILE"

echo "Adding answer files to context:" >>"$LOG_FILE"
for afile in "$ANSWER_DIR"/*; do
  if [ -f "$afile" ] && [ "$afile" != "$ANSWER_FILE" ]; then
    CONTEXT_ARGS+=("--context-file" "$afile")
    echo "  Added answer file to context: $afile" >>"$LOG_FILE"
  elif [ "$afile" = "$ANSWER_FILE" ]; then
    echo "  Skipping target answer file: $afile" >>"$LOG_FILE"
  fi
done

echo "All context arguments (${#CONTEXT_ARGS[@]} args):" >>"$LOG_FILE"
for ((i = 0; i < ${#CONTEXT_ARGS[@]}; i += 2)); do
  echo "  ${CONTEXT_ARGS[i]} ${CONTEXT_ARGS[i + 1]}" >>"$LOG_FILE"
done

echo "Executing Cody chat command..." >>"$LOG_FILE"

echo "Command to execute: cody chat --silent --model \"$MODEL\" [context_args] --message \"[Prompt from $PROMPT_FILE]\"" >>"$LOG_FILE"

if TERM=dumb cody chat --model "$MODEL" "${CONTEXT_ARGS[@]}" --message "$PROMPT" >"$ANSWER_FILE" 2>>"$LOG_FILE"; then
  echo "Cody chat command executed successfully." >>"$LOG_FILE"
  echo "Câu trả lời đã được lưu vào: $ANSWER_FILE" >>"$LOG_FILE"
  echo "Output saved to: $ANSWER_FILE" >>"$LOG_FILE"
else
  echo "ERROR: Cody chat command failed. Check $LOG_FILE for details." | tee -a "$LOG_FILE"
  exit 1
fi

echo "Script finished at: $(date)" >>"$LOG_FILE"
echo "==================================================" >>"$LOG_FILE"
echo "" >>"$LOG_FILE"
