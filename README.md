This action installs and configures [Aliyun command line tool](https://github.com/aliyun/aliyun-cli) for use in your GitHub
Action steps.

## Usage

```yaml
steps:
- uses: actions/checkout@v1
- uses: aliyun/aliyun-cli-action@v1.0.0
  with:
    access-key-id: ${{ secrets.ALIYUN_ACCESS_KEY_ID }}
    access-key-secret: ${{ secrets.ALIYUN_ACCESS_KEY_SECRET }}
    region: ${{ secrets.ALIYUN_REGION }}
- run: aliyun oss cp ./dir oss://backet/path -r -u
```

## License

MIT
