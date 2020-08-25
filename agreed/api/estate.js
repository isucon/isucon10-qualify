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

  // `GET: /api/estate/search/condition`
  {
    request: {
      path: `${PATH}/search/condition`,
      method: 'GET',
      query: {},
      values: {}
    },
    response: {
      body: {
        doorWidth: '{:doorWidth}',
        doorHeight: '{:doorHeight}',
        rent: '{:rent}'
      },
      schema: {
        type: 'object',
        properties: {
          doorWidth: {
            type: 'object',
            properties: {
              prefix: 'string',
              suffix: 'string',
              ranges: {
                type: 'array',
                items: {
                  type: 'object',
                  properties: {
                    id: 'number',
                    min: 'number',
                    max: 'number'
                  }
                }
              }
            }
          },
          doorHeight: {
            type: 'object',
            properties: {
              prefix: 'string',
              suffix: 'string',
              ranges: {
                type: 'array',
                items: {
                  type: 'object',
                  properties: {
                    id: 'number',
                    min: 'number',
                    max: 'number'
                  }
                }
              }
            }
          },
          rent: {
            type: 'object',
            properties: {
              prefix: 'string',
              suffix: 'string',
              ranges: {
                type: 'array',
                items: {
                  type: 'object',
                  properties: {
                    id: 'number',
                    min: 'number',
                    max: 'number'
                  }
                }
              }
            }
          }
        }
      },
      values: {
        doorWidth: {
          prefix: '',
          suffix: 'cm',
          ranges: [
            {
              id: 0,
              min: -1,
              max: 80
            },
            {
              id: 1,
              min: 81,
              max: 110
            },
            {
              id: 2,
              min: 111,
              max: 150
            },
            {
              id: 3,
              min: 151,
              max: -1
            }
          ]
        },
        doorHeight: {
          prefix: '',
          suffix: 'cm',
          ranges: [
            {
              id: 0,
              min: -1,
              max: 80
            },
            {
              id: 1,
              min: 81,
              max: 110
            },
            {
              id: 2,
              min: 111,
              max: 150
            },
            {
              id: 3,
              min: 151,
              max: -1
            }
          ]
        },
        rent: {
          prefix: '',
          suffix: '円',
          ranges: [
            {
              id: 0,
              min: -1,
              max: 50000
            },
            {
              id: 1,
              min: 50001,
              max: 100000
            },
            {
              id: 2,
              min: 100001,
              max: 150000
            },
            {
              id: 3,
              min: 150001,
              max: -1
            }
          ]
        }
      }
    }
  },

  // GET: /api/estate/low_priced
  {
    request: {
      path: `${PATH}/low_priced`,
      method: 'GET',
      body: {},
      values: {}
    },
    response: {
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
                thumbnail: 'string',
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
            thumbnail: '/images/estate/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
            name: 'イスイスレジデンス南タワー',
            description: 'ビル群の中に佇む最高のお部屋、さらなるイスの高みへ',
            address: '東京都千代田区丸の内1丁目9-2',
            latitude: 35.678637,
            longitude: 139.767375,
            doorHeight: 230,
            doorWidth: 120,
            rent: 2500000,
            features: '駅直結,バストイレ別'
          },
          {
            id: 5,
            thumbnail: '/images/estate/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg',
            name: '四丼往親空中イスコビル',
            description: '一階が金融機関になっております！',
            address: '東京都中央区京橋1丁目6-1',
            latitude: 35.678617,
            longitude: 139.767345,
            doorHeight: 220,
            doorWidth: 150,
            rent: 2000000,
            features: '音響攻撃あり,バストイレ別'
          }
        ]
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
        features: '{:features}',
        page: '{:page}',
        perPage: '{:perPage}'
      },
      values: {
        rentRangeId: 2,
        doorHeightRangeId: 3,
        doorWidthRangeId: 2,
        features: 'バストイレ別,DIY可',
        page: 0,
        perPage: 20
      }
    },
    response: {
      body: {
        count: '{:count}',
        estates: '{:estates}'
      },
      schema: {
        type: 'object',
        properties: {
          count: 'number',
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
        count: 2000,
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
            popularity: 10000
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
            popularity: 10000
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
      status: 200,
      body: 'OK',
      schema: {
        type: 'string'
      },
      values: {}
    }
  }
]
