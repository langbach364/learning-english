function parseDefinitions(content) {
  const lines = content.split('\n');
  const definitions = {};
  let currentWord = '';
  let currentDefinition = {};

  lines.forEach(line => {
    line = line.trim();
    if (line.startsWith('- Từ:')) {
      currentWord = line.split(':')[1].trim();
      definitions[currentWord] = [];
    } else if (line.startsWith('*')) {
      if (Object.keys(currentDefinition).length > 0) {
        definitions[currentWord].push(currentDefinition);
      }
      currentDefinition = {};
      const [type, content] = line.substring(1).split(':').map(s => s.trim());
      currentDefinition['Từ loại'] = type;
      currentDefinition['Định nghĩa'] = content;
      currentDefinition['Ví dụ'] = [];
    } else if (line.startsWith('+')) {
      const example = line.substring(1).split(':')[1].trim();
      currentDefinition['Ví dụ'].push(example);
    }
  });

  if (Object.keys(currentDefinition).length > 0) {
    definitions[currentWord].push(currentDefinition);
  }

  return definitions;
}

window.parseDefinitions = parseDefinitions;
