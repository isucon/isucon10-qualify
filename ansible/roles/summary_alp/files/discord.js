"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const child_process_1 = require("child_process");
const fs_1 = __importDefault(require("fs"));
const alpResultFilename = "alp_result.txt";
const regexp = "/api/user/[a-zA-Z0-9]+/icon,/api/livestream/[0-9]+/statistics,/api/livestream/[0-9]+/livecomment,/api/livestream/[0-9]+/reaction,/api/livestream/[0-9]+/report,/api/livestream/[0-9]+/exit,/api/livestream/[0-9]+/enter,/api/livestream/[0-9]+/moderate,/api/livestream/[0-9]+/ngwords,/api/user/[a-zA-Z0-9]+/statistics,/api/user/[a-zA-Z0-9]+/theme,/assets/.*,/api/upload/.*";
const webhookUrl = "https://discord.com/api/webhooks/1173915786243477574/XgRRPJai3X5XcTgxPr-0AyzxsaWWdE3gB3NBBz_ZnmmmZjBS4gI6JD_E65G2bHlW3_Da";
const alpCommand = `sudo cat /var/log/nginx/access.log | alp json -m "${regexp}"`;
const curlCommand = `curl -H "Content-type: multipart/form-data" -X POST -F "file=@${alpResultFilename}" '${webhookUrl}'`;
/*
 *  Discordにalpの結果をテキストファイルにして送信する
 */
const main = async () => {
    const alpResult = (0, child_process_1.execSync)(alpCommand).toString();
    fs_1.default.writeFileSync(alpResultFilename, alpResult);
    const res = (0, child_process_1.execSync)(curlCommand).toString();
    console.log(res);
};
main();
//# sourceMappingURL=main.js.map