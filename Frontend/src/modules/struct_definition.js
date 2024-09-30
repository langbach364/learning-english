function parseDefinitions(content) {
  const lines = content.split('\n');
  const result = { sentence: {}, words: {} };
  let currentSection = '';
  let currentWord = '';
  let currentType = '';
  let inWordSection = false;

  lines.forEach(line => {
    line = line.trim();
    if (line.startsWith('- Câu:')) {
      inWordSection = false;
      result.sentence['Câu gốc'] = line.substring(7).trim();
    } else if (line.startsWith('[Từ loại]:')) {
      inWordSection = true;
    } else if (inWordSection && line.startsWith('- Từ:')) {
      currentWord = line.substring(6).trim();
    } else if (inWordSection && line.startsWith('*')) {
      currentType = line.substring(1).split(':')[0].trim();
      let definition = line.includes(':') ? line.split(':')[1].trim() : '';
      if (!result.words[currentType]) {
        result.words[currentType] = {};
      }
      if (!result.words[currentType][currentWord]) {
        result.words[currentType][currentWord] = [];
      }
      result.words[currentType][currentWord].push({
        definition: definition,
        examples: []
      });
    } else if (inWordSection && line.startsWith('+')) {
      let lastDefinition = result.words[currentType][currentWord][result.words[currentType][currentWord].length - 1];
      lastDefinition.examples.push(line.substring(1).trim());
    } else if (!inWordSection) {
      if (line.startsWith('*')) {
        currentSection = line.substring(1).split(':')[0].trim();
        result.sentence[currentSection] = line.includes(':') ? line.split(':')[1].trim() : [];
      } else if (line.startsWith('+') || line.startsWith('-')) {
        if (Array.isArray(result.sentence[currentSection])) {
          result.sentence[currentSection].push(line.substring(1).trim());
        } else {
          result.sentence[currentSection] = [line.substring(1).trim()];
        }
      } else if (line && currentSection) {
        if (typeof result.sentence[currentSection] === 'string') {
          result.sentence[currentSection] += ' ' + line;
        } else if (Array.isArray(result.sentence[currentSection])) {
          result.sentence[currentSection][result.sentence[currentSection].length - 1] += ' ' + line;
        }
      }
    }
  });

  return result;
}

window.parseDefinitions = parseDefinitions;
