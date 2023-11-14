#!/usr/bin/env node

import { writeFileSync } from 'fs';
import { program, OptionValues } from 'commander';

const generateCode = (opt: OptionValues) => {

  const dataTypes = opt.types.split(',');

  let code = `import { MongoLocator } from "privacy-pal";

export default function handleAccess(dataSubjectId: string, locator: MongoLocator, obj: any): Record<string, any> {

  switch (locator.dataType) {
`;

  dataTypes.forEach((dataType: string) => {
    code += `
    case "${dataType}":
      return handleAccess${capitalize(dataType)}(dataSubjectId, locator, obj)
`;
  });

  code += `
  }
  return {};
}`;

  dataTypes.forEach((dataType: string) => {
    code += `

function handleAccess${capitalize(dataType)}(dataSubjectId: string, locator: MongoLocator, obj: any): Record<string, any> {
  return {};
}
`;
  });

  writeFileSync('./generated.ts', code);

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