#!/usr/bin/env node

import { writeFileSync } from 'fs';
import { format } from 'prettier';
import { program, OptionValues } from 'commander';

const generateCode = (opt: OptionValues) => {

  const dataTypes = opt.types.split(',');

  const cases = dataTypes.map((dataType: string) => (`
  case "${dataType}":
    return handleAccess${capitalize(dataType)}(dataSubjectId, locator, obj)
`)).join('');

  const funcs = dataTypes.map((dataType: string) => (`function handleAccess${capitalize(dataType)}(dataSubjectId: string, locator: MongoLocator, obj: any): Record<string, any> {
  return {};
}

`)).join('');

  const code = `import { MongoLocator } from "privacy-pal";

export default function handleAccess(dataSubjectId: string, locator: MongoLocator, obj: any): Record<string, any> {

  switch (locator.dataType) {
    ${cases}
  }

  return {};
}

${funcs}
`;

  format(code, { parser: 'typescript' })
  .then((formattedCode) => writeFileSync(opt.output, formattedCode))
  .catch((err) => console.error(err));

}

const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

const main = () => {
  program
    .option('-t, --types <type>', 'comma-separated list of data types')
    .option('-o, --output <file>', 'output file', './generated.ts');

  program.parse(process.argv);
  const options = program.opts();

  generateCode(options);
}

main();