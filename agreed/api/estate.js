const PATH = '/api/estate'

module.exports = [
  // `GET: /api/estate/:id`
  {
    request: {
      path: `${PATH}/:id`,
      method: 'GET',
      body: {},
      values: {}
    },
    response: {
      headers: {
        'Access-Control-Allow-Origin': '*'
      },
      status: 200,
      body: {
        id: '{:id}',
        thumbnail: '{:thumbnail}',
        name: '{:name}',
        description: '{:description}',
        address: '{:address}',
        latitude: '{:latitude}',
        longitude: '{:longitude}',
        doorHeight: '{:doorHeight}',
        doorWidth: '{:doorWidth}',
        rent: '{:rent}',
        features: '{:features}'
      },
      schema: {
        type: 'object',
        properties: {
          id: 'number',
          thumbnail: 'string',
          name: 'string',
          description: 'string',
          address: 'string',
          latitude: 'number',
          longitude: 'number',
          doorHeight: 'number',
          doorWidth: 'number',
          rent: 'number',
          features: 'string'
        }
      },
      values: {
        id: 1,
        thumbnail: '/images/estate/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
        name: '公園前派出所',
        description: '両津勘吉',
        address: '東京都葛飾区亀有',
        latitude: 37,
        longitude: 137,
        doorHeight: 200,
        doorWidth: 100,
        rent: 40000,
        features: 'バストイレ別,DIY可'
      }
    }
  },

  // `GET: /api/estate/search`
  {
    request: {
      path: `${PATH}/search`,
      method: 'GET',
      query: {
        rentRangeId: '{:priceRangeId}',
        doorHeightRangeId: '{:doorHeightId}',
        doorWidthRangeId: '{:doorWidthId}',
        features: '{:features}'
      },
      values: {
        rentRangeId: 2,
        doorHeightRangeId: 3,
        doorWidthRangeId: 2,
        features: 'バストイレ別,DIY可'
      }
    },
    response: {
      headers: {
        'Access-Control-Allow-Origin': '*'
      },
      body: {
        estates: '{:estates}'
      },
      schema: {
        type: 'object',
        items: {
          estates: {
            type: 'array',
            items: {
              type: 'object',
              properties: {
                id: 'number',
                thumbnail: 'string',
                name: 'string',
                description: 'string',
                address: 'string',
                latitude: 'number',
                longitude: 'number',
                doorHeight: 'number',
                doorWidth: 'number',
                rent: 'number',
                features: 'string'
              }
            }
          }
        }
      },
      values: {
        estates: [
          {
            id: 1,
            thumbnail: '/images/estate/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg',
            name: 'イスイスレジデンス南タワー',
            description: 'ビル群の中に佇む最高のお部屋、さらなるイスの高みへ',
            address: '東京都千代田区丸の内1丁目9-2',
            latitude: 35,
            longitude: 137,
            doorHeight: 230,
            doorWidth: 120,
            rent: 2500000,
            features: '駅直結,バストイレ別',
            view_count: 10000
          },
          {
            id: 5,
            thumbnail: '/images/estate/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
            name: '四丼往親空中イスコビル',
            description: '一階が金融機関になっております！',
            address: '東京都中央区京橋1丁目6-1',
            latitude: 35,
            longitude: 135,
            doorHeight: 220,
            doorWidth: 150,
            rent: 2000000,
            features: '音響攻撃あり,バストイレ別'
          }
        ]
      }
    }
  },
  {
    request: {
      path: `${PATH}/nazotte`,
      method: 'POST',
      body: {
        coordinates: '{:coordinates}'
      },
      values: {
        coordinates: [
          {
            latitude: 36.5,
            longitude: 137.5
          },
          {
            latitude: 36.5,
            longitude: 138.5
          },
          {
            latitude: 37.5,
            longitude: 138.5
          },
          {
            latitude: 37.5,
            longitude: 137.5
          },
          {
            latitude: 36.5,
            longitude: 137.5
          }
        ]
      }
    },
    response: {
      headers: {
        'Access-Control-Allow-Origin': '*'
      },
      body: {
        estates: '{:estates}'
      },
      schema: {
        type: 'object',
        items: {
          estates: {
            type: 'array',
            items: {
              type: 'object',
              properties: {
                id: 'number',
                thumbnail: 'string',
                name: 'string',
                description: 'string',
                address: 'string',
                latitude: 'number',
                longitude: 'number',
                doorHeight: 'number',
                doorWidth: 'number',
                rent: 'number',
                features: 'string'
              }
            }
          }
        }
      },
      values: {
        estates: [
          {
            id: 1,
            thumbnail: '/images/estate/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
            name: 'イスイスレジデンス南タワー',
            description: 'ビル群の中に佇む最高のお部屋、さらなるイスの高みへ',
            address: '東京都千代田区丸の内1丁目9-2',
            latitude: 35,
            longitude: 135,
            doorHeight: 230,
            doorWidth: 120,
            rent: 2500000,
            features: '駅直結,バストイレ別',
            view_count: 10000
          },
          {
            id: 5,
            thumbnail: '/images/estate/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg',
            name: '四丼往親空中イスコビル',
            description: '一階が金融機関になっております！',
            address: '東京都中央区京橋1丁目6-1',
            latitude: 35,
            longitude: 137,
            doorHeight: 220,
            doorWidth: 150,
            rent: 2000000,
            features: '音響攻撃あり,バストイレ別'
          }
        ]
      }
    }
  },

  // `POST: /api/estate/req_doc/:id`
  {
    request: {
      path: `${PATH}/req_doc/:id`,
      method: 'POST',
      body: {
        email: '{:email}'
      },
      values: {
        id: 10,
        email: 'isuumo@example.com'
      }
    },
    response: {
      headers: {
        'Access-Control-Allow-Origin': '*'
      },
      status: 200,
      body: 'OK',
      schema: {
        type: 'string'
      },
      values: {}
    }
  }
]
