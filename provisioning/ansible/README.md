# isucon10-provisioning

ansible 2.9.13で動作確認しています

## playbooks
- benchmarker.yaml
  - ベンチマーカーがセットアップされます
- competitor.yaml
  - 競技者に提供された各種言語実装がセットアップされます
- allinone.yaml
  - 各種言語実装に加えてベンチマーカーのセットアップもされています

## Vagrantを利用して，環境をセットアップする

本Vagrantファイルは1台構成で，allinone.yamlを実行した結果を提供しています
別の，設定を利用したい場合は,
- Vagrant ファイルの書き換え
      ansible.playbook = "allinone.yaml"
- inventory/hostsの書き換え
を行ってから，下記の操作を行ってください．

### 初回構築
下記コマンドで，VMを一度破棄して，新しく作り直します

make vagrant/init

### ファイル初期化/再構築
ansibleの実行中に，通信環境エラーなどが起きた場合，下記のコマンドで再実行できます
vagrant provision


## サーバーへのprovisionning
- inventory/hostsを書き換えて，ansible playbookを実行してください．
