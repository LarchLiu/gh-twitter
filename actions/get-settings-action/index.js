import * as core from '@actions/core';
import path from 'path';
import fs from 'fs';

const run = () => {
  const settings_path = core.getInput('settings_path', { required: true });
  const filePath = path.resolve(settings_path)
  if (fs.existsSync(filePath)) {
    var fileStr = fs.readFileSync(filePath, 'utf-8');
    var jsonObj = JSON.parse(fileStr);
    var keys = Object.keys(jsonObj)
    for (var i in keys) {
      core.setOutput(keys[i], jsonObj[keys[i]]);
    }
    core.setOutput("SettingsPath", filePath)
  } else {
    core.setFailed('file not exist');
  }
};

if (process.env.NODE_ENV !== 'test') {
  run();
}

export default run;