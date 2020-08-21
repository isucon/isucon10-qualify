const PATH = '/api'

module.exports = [
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
            thumbnail: '/images/chair/1501E5C34A2B8EE645480ED1CC6442CD5929FE7616E20513574628096163DF0C.jpg',
            name: 'isuu megro',
            description: 'ビル群の中に佇む最高のお部屋、さらなるイスの高みへ',
            address: '東京都品川区上大崎2丁目13-30',
            latitude: 35.678637,
            longitude: 139.767375,
            doorHeight: 230,
            doorWidth: 120,
            rent: 2500000,
            features: '駅直結,バストイレ別'
          },
          {
            id: 5,
            thumbnail: '/images/chair/3E880A828B1DBFACB42209724583B56EF28466E45E2BF3704475EA02B19BDBFC.jpg',
            name: 'イスリック銀座7丁目ビル',
            description: '一階が車のディーラーになっております！',
            address: '東京都中央区銀座7-3-5',
            latitude: 35.678617,
            longitude: 139.767345,
            doorHeight: 220,
            doorWidth: 150,
            rent: 2000000,
            features: '便利な好立地,バストイレ別'
          }
        ]
      }
    }
  }
]
