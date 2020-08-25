const PATH = '/api/chair'

module.exports = [
  // GET: /api/chair/:id
  {
    request: {
      path: `${PATH}/:id`,
      method: 'GET',
      query: {},
      values: {
        id: 10
      }
    },
    response: {
      status: 200,
      body: {
        id: '{:id}',
        name: '{:name}',
        description: '{:description}',
        thumbnail: '{:thumbnail}',
        price: '{:price}',
        height: '{:height}',
        width: '{:width}',
        depth: '{:depth}',
        color: '{:color}',
        features: '{:features}',
        kind: '{:kind}'
      },
      schema: {
        type: 'object',
        properties: {
          id: 'number',
          name: 'string',
          description: 'string',
          thumbnail: 'string',
          price: 'number',
          height: 'number',
          width: 'number',
          depth: 'number',
          color: 'string',
          features: 'string',
          kind: 'string'
        }
      },
      values: {
        id: 10,
        name: 'すごいイス',
        description: 'すごいネコはいます。',
        thumbnail: '/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
        price: 10000,
        height: 100,
        width: 50,
        depth: 60,
        color: '黒',
        features: 'リクライニング,キャスター付き,肘掛け',
        kind: 'エルゴノミクス'
      }
    }
  },

  // `GET: /api/chair/search/condition`
  {
    request: {
      path: `${PATH}/search/condition`,
      method: 'GET',
      query: {},
      values: {}
    },
    response: {
      body: {
        price: '{:price}',
        height: '{:height}',
        width: '{:width}',
        depth: '{:depth}'
      },
      schema: {
        type: 'object',
        properties: {
          price: {
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
          width: {
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
          height: {
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
          depth: {
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
        price: {
          prefix: '',
          suffix: '円',
          ranges: [
            {
              id: 0,
              min: -1,
              max: 3000
            },
            {
              id: 1,
              min: 3001,
              max: 6000
            },
            {
              id: 2,
              min: 6001,
              max: 9000
            },
            {
              id: 3,
              min: 9001,
              max: 12000
            },
            {
              id: 4,
              min: 12001,
              max: 15000
            },
            {
              id: 5,
              min: 15001,
              max: -1
            }
          ]
        },
        width: {
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
        height: {
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
        depth: {
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
        }
      }
    }
  },

  // POST: /api/chair/low_priced
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
        chairs: '{:chairs}'
      },
      schema: {
        type: 'object',
        properties: {
          id: 'number',
          name: 'string',
          description: 'string',
          thumbnail: 'string',
          price: 'number',
          height: 'number',
          width: 'number',
          depth: 'number',
          color: 'string',
          features: 'string',
          kind: 'string'
        }
      },
      values: {
        chairs: [
          {
            id: 10,
            name: 'スモスモチェアー',
            description: 'スモスモハウスにぴったりの素敵なイスです',
            thumbnail: '/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
            price: 10000,
            height: 100,
            width: 50,
            depth: 60,
            color: '緑',
            features: 'リクライニング,キャスター付き,肘掛け',
            kind: 'エルゴノミクス'
          },
          {
            id: 13,
            name: '王様のイス',
            description: 'どうぶつの森からきました',
            thumbnail: '/images/chair/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
            price: 100000,
            height: 100,
            width: 50,
            depth: 60,
            color: '黄',
            features: 'リクライニング,キャスター付き,肘掛け',
            kind: 'エルゴノミクス'
          }
        ]
      }
    }
  },

  // GET: /api/chair/search
  {
    request: {
      path: `${PATH}/search`,
      method: 'GET',
      query: {
        priceRangeId: '{:priceRangeId}',
        heightRangeId: '{:height}',
        widthRangeId: '{:width}',
        depthRangeId: '{:depth}',
        color: '{:color}',
        features: '{:features}',
        kind: '{:kind}',
        page: '{:page}',
        perPage: '{:perPage}'
      },
      values: {
        priceRangeId: 2,
        heightRangeId: 3,
        widthRangeId: 2,
        depthRangeId: 1,
        color: '黒',
        features: 'リクライニング,肘掛け',
        kind: 'エルゴノミクス',
        page: 0,
        perPage: 20
      }
    },
    response: {
      body: {
        count: '{:count}',
        chairs: '{:chairs}'
      },
      values: {
        count: 2000,
        chairs: [
          {
            id: 1,
            name: 'すごいイス',
            description: 'すごいネコはいます。',
            thumbnail: '/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
            price: 10000,
            height: 100,
            width: 50,
            depth: 60,
            color: '黒',
            features: 'リクライニング,キャスター付き,肘掛け',
            kind: 'エルゴノミクス'
          },
          {
            id: 11,
            name: 'ボロいイス',
            description: 'ボロい釣り竿的なsomething。',
            thumbnail: '/images/chair/9120C2E3CAF5CD376C1B14899C2FD31438A839D1F6B6F8A52091392E0B9168FC.jpg',
            price: 12000,
            height: 80,
            width: 45,
            depth: 70,
            color: '黒',
            features: '肘掛け',
            kind: 'エルゴノミクス'
          },
          {
            id: 12,
            name: 'ふつうのハンモック',
            description: '老後はハンモックで遊びたい。',
            thumbnail: '/images/chair/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg',
            price: 12000,
            height: 50,
            width: 120,
            depth: 70,
            color: '白',
            features: 'リクライニング',
            kind: 'ハンモック'
          }
        ]
      }
    }
  },

  // POST: /api/chair/buy/:id
  {
    request: {
      path: `${PATH}/buy/:id`,
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
      body: 'OK',
      schema: {
        type: 'string'
      },
      values: {}
    }
  }
]
