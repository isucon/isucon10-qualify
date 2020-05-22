const PATH = '/api'

module.exports = [
  // GET: /api/recommended_estate
  {
    request: {
      path: `${PATH}/recommended_estate`,
      method: 'GET',
      body: {},
      values: {}
    },
    response: {
      headers: {},
      status: 200,
      body: {
        estates: '{:estates}'
      },
      schema: {
        type: 'object',
        properties: {
          estates: {
            type: 'array',
            items: {
              type: 'object',
              properties: {
                id: 'number',
                thumbnails: 'array',
                name: 'string',
                description: 'string',
                address: 'string',
                latitude: 'number',
                longitude: 'number',
                doorHeight: 'number',
                doorWidth: 'number',
                rent: 'number',
                features: 'array'
              }
            }
          }
        }
      },
      values: {
        estates: [
          {
            id: 1,
            thumbnails: [
              '/assets/images/estate/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
              '/assets/images/estate/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
              '/assets/images/estate/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg'
            ],
            name: 'イスイスレジデンス南タワー',
            description: 'ビル群の中に佇む最高のお部屋、さらなるイスの高みへ',
            address: '東京都千代田区丸の内1丁目9-2',
            latitude: 35.678637,
            longitude: 139.767375,
            doorHeight: 230,
            doorWidth: 120,
            rent: 2500000,
            features: [
              '駅直結',
              'バストイレ別'
            ]
          },
          {
            id: 5,
            thumbnails: [
              '/assets/images/estate/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
              '/assets/images/estate/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
              '/assets/images/estate/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg'
            ],
            name: '四丼往親空中イスコビル',
            description: '一階が金融機関になっております！',
            address: '東京都中央区京橋1丁目6-1',
            latitude: 35.678617,
            longitude: 139.767345,
            doorHeight: 220,
            doorWidth: 150,
            rent: 2000000,
            features: [
              '音響攻撃あり',
              'バストイレ別'
            ]
          }
        ]
      }
    }
  },

  // POST: /api/recommended_chair
  {
    request: {
      path: `${PATH}/recommended_chair`,
      method: 'GET',
      body: {},
      values: {}
    },
    response: {
      headers: {},
      status: 200,
      body: {
        chairs: '{:chairs}'
      },
      schema: {
        type: 'object',
        properties: {
          id: 'number',
          name: 'string',
          description: 'string',
          thumbnails: {
            type: 'array',
            items: 'string'
          },
          price: 'number',
          height: 'number',
          width: 'number',
          depth: 'number',
          color: 'string',
          features: {
            type: 'array',
            items: 'string'
          },
          kind: 'string'
        }
      },
      values: {
        chairs: [
          {
            id: 10,
            name: 'スモスモチェアー',
            description: 'スモスモハウスにぴったりの素敵なイスです',
            thumbnails: [
              '/assets/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
              '/assets/images/chair/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
              '/assets/images/chair/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg'
            ],
            price: 10000,
            height: 100,
            width: 50,
            depth: 60,
            color: '緑',
            features: [
              'リクライニング',
              'キャスター付き',
              '肘掛け'
            ],
            kind: 'エルゴノミクス'
          },
          {
            id: 13,
            name: '王様のイス',
            description: 'どうぶつの森からきました',
            thumbnails: [
              '/assets/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
              '/assets/images/chair/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
              '/assets/images/chair/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg'
            ],
            price: 100000,
            height: 100,
            width: 50,
            depth: 60,
            color: '黄',
            features: [
              'リクライニング',
              'キャスター付き',
              '肘掛け'
            ],
            kind: 'エルゴノミクス'
          }
        ]
      }
    }
  },

  // POST: /api/recommended_estate/:chairId
  {
    request: {
      path: `${PATH}/recommended_estate/:chairId`,
      method: 'GET',
      body: {},
      values: {
        chairId: 10
      }
    },
    response: {
      headers: {},
      status: 200,
      body: {
        estates: '{:estates}'
      },
      schema: {
        type: 'object',
        properties: {
          estates: {
            type: 'array',
            items: {
              type: 'object',
              estates: {
                id: 'number',
                thumbnails: 'array',
                name: 'string',
                description: 'string',
                address: 'string',
                latitude: 'number',
                longitude: 'number',
                doorHeight: 'number',
                doorWidth: 'number',
                rent: 'number',
                features: 'array'
              }
            }
          }
        }
      },
      values: {
        estates: [
          {
            id: 1,
            thumbnails: [
              '/assets/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
              '/assets/images/chair/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
              '/assets/images/chair/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg'
            ],
            name: 'isuu megro',
            description: 'ビル群の中に佇む最高のお部屋、さらなるイスの高みへ',
            address: '東京都品川区上大崎2丁目13-30',
            latitude: 35.678637,
            longitude: 139.767375,
            doorHeight: 230,
            doorWidth: 120,
            rent: 2500000,
            features: [
              '駅直結',
              'バストイレ別'
            ]
          },
          {
            id: 5,
            thumbnails: [
              '/assets/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
              '/assets/images/chair/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
              '/assets/images/chair/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg'
            ],
            name: 'イスリック銀座7丁目ビル',
            description: '一階が車のディーラーになっております！',
            address: '東京都中央区銀座7-3-5',
            latitude: 35.678617,
            longitude: 139.767345,
            doorHeight: 220,
            doorWidth: 150,
            rent: 2000000,
            features: [
              '便利な好立地',
              'バストイレ別'
            ]
          }
        ]
      }
    }
  }
]
