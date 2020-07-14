import { useEffect, useState } from 'react'
import Link from 'next/link'
import {
  Paper,
  Container,
  Box,
  Button,
  CircularProgress,
  Card,
  CardContent,
  CardMedia,
  CardActionArea
} from '@material-ui/core'
import { Pagination } from '@material-ui/lab'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import { Loading } from '../../../components/Loading'
import { RangeForm } from '../../../components/RangeForm'
import { RadioButtonForm } from '../../../components/RadioButtonForm'
import { CheckboxForm } from '../../../components/CheckboxForm'

import type { FC } from 'react'
import type { ChairRangeMap, ChairSearchCondition, ChairSearchResponse } from '@types'

const ESTATE_COUNTS_PER_PAGE = 20

interface ChairSearchProps {
  chairRangeMap: ChairRangeMap
}

const useChairSearchStyles = makeStyles(theme =>
  createStyles({
    page: {
      margin: theme.spacing(2),
      padding: theme.spacing(4)
    },
    search: {
      display: 'flex',
      flexDirection: 'column',
      marginTop: theme.spacing(4),
      marginBottom: theme.spacing(4),
      '&>*': {
        margin: theme.spacing(1)
      }
    },
    row: {
      '&>*': {
        margin: theme.spacing(2)
      }
    },
    card: {
      width: '100%',
      height: 270,
      marginTop: theme.spacing(2),
      marginBottom: theme.spacing(2)
    },
    cardActionArea: {
      display: 'flex',
      alignItems: 'flex-start',
      justifyContent: 'flex-start'
    },
    cardMedia: {
      width: 360,
      height: 270
    },
    cardContent: {
      marginLeft: theme.spacing(1)
    }
  })
)

const COLOR_LIST = [
  '黒',
  '白',
  '赤',
  '青',
  '緑',
  '黄',
  '紫',
  'ピンク',
  'オレンジ',
  '水色',
  'ネイビー',
  'ベージュ'
]

const FEATURE_LIST = [
  '折りたたみ可',
  '肘掛け',
  'キャスター',
  'リクライニング',
  '高さ調節可',
  'フットレスト'
]

const KIND_LIST = [
  'ゲーミングチェア',
  '座椅子',
  'エルゴノミクス',
  'ハンモック'
]

const ChairSearch: FC<ChairSearchProps> = ({ chairRangeMap }) => {
  const classes = useChairSearchStyles()

  const [priceRangeId, setPriceRangeId] = useState('')
  const [heightRangeId, setHeightRangeId] = useState('')
  const [widthRangeId, setWidthRangeId] = useState('')
  const [depthRangeId, setDepthRangeId] = useState('')
  const [color, setColor] = useState('')
  const [kind, setKind] = useState('')
  const [features, setFeatures] = useState<boolean[]>(new Array(FEATURE_LIST.length).fill(false))
  const [chairSearchCondition, setChairSearchCondition] = useState<ChairSearchCondition | null>(null)
  const [searchResult, setSearchResult] = useState<ChairSearchResponse | null>(null)
  const [page, setPage] = useState<number>(0)

  const onSearch = () => {
    const selectedFeatures = FEATURE_LIST.filter((_, i) => features[i])
    const condition: ChairSearchCondition = {
      priceRangeId,
      heightRangeId,
      widthRangeId,
      depthRangeId,
      color,
      kind,
      features: selectedFeatures.length > 0 ? selectedFeatures.join(',') : '',
      page: 0,
      perPage: ESTATE_COUNTS_PER_PAGE
    }
    setChairSearchCondition(condition)

    const params = new URLSearchParams()
    for (const [key, value] of Object.entries(condition)) {
      params.append(key, value.toString())
    }
    fetch(`/api/chair/search?${params.toString()}`, { mode: 'cors' })
      .then(async response => await response.json())
      .then(result => {
        setSearchResult(result as ChairSearchResponse)
        setPage(0)
      })
      .catch(console.error)
  }

  return (
    <>
      <Paper className={classes.page}>
        <Container maxWidth='md'>
          <Box width={1} className={classes.search}>
            <RangeForm
              name='イスの高さ'
              value={heightRangeId}
              rangeList={chairRangeMap.height}
              onChange={(_, value) => { setHeightRangeId(value) }}
            />

            <RangeForm
              name='イスの横幅'
              value={widthRangeId}
              rangeList={chairRangeMap.width}
              onChange={(_, value) => { setWidthRangeId(value) }}
            />

            <RangeForm
              name='イスの奥行き'
              value={depthRangeId}
              rangeList={chairRangeMap.depth}
              onChange={(_, value) => { setDepthRangeId(value) }}
            />

            <RangeForm
              name='価格'
              value={priceRangeId}
              rangeList={chairRangeMap.price}
              onChange={(_, value) => { setPriceRangeId(value) }}
            />

            <RadioButtonForm
              name='色'
              value={color}
              items={COLOR_LIST}
              onChange={(_, value) => { setColor(value) }}
            />

            <RadioButtonForm
              name='種類'
              value={kind}
              items={KIND_LIST}
              onChange={(_, value) => { setKind(value) }}
            />

            <CheckboxForm
              name='特徴'
              checkList={features}
              selectList={FEATURE_LIST}
              onChange={(_, checked, key) => {
                setFeatures(features.map((feature, i) => key === i ? checked : feature))
              }}
            />

            <Button
              onClick={onSearch}
              disabled={
                heightRangeId === '' &&
                widthRangeId === '' &&
                depthRangeId === '' &&
                priceRangeId === '' &&
                color === '' &&
                kind === '' &&
                !features.some(feature => feature)
              }
            >
              Search
            </Button>
          </Box>
        </Container>
      </Paper>

      {chairSearchCondition ? (
        <Paper className={classes.page}>
          <Container maxWidth='md'>
            <Box width={1} className={classes.search} alignItems='center'>
              {searchResult ? (
                <>
                  <Pagination
                    count={Math.ceil(searchResult.count / ESTATE_COUNTS_PER_PAGE)}
                    page={page + 1}
                    onChange={(_, page) => {
                      setPage(page - 1)
                      if (!chairSearchCondition) return
                      const condition = { ...chairSearchCondition, page: page - 1 }
                      setChairSearchCondition(condition)
                      setSearchResult(null)
                      const params = new URLSearchParams()
                      for (const [key, value] of Object.entries(condition)) {
                        params.append(key, value.toString())
                      }
                      fetch(`/api/chair/search?${params.toString()}`, { mode: 'cors' })
                        .then(async response => await response.json())
                        .then(result => {
                          setSearchResult(result as ChairSearchResponse)
                        })
                        .catch(console.error)
                    }}
                  />
                  {
                    searchResult.chairs.map((chair) => (
                      <Link key={chair.id} href={`/chair/detail?id=${chair.id}`}>
                        <Card className={classes.card}>
                          <CardActionArea className={classes.cardActionArea}>
                            <CardMedia image={chair.thumbnail} className={classes.cardMedia} />
                            <CardContent className={classes.cardContent}>
                              <h2>{chair.name}</h2>
                              <p>価格: {chair.price}円</p>
                              <p>詳細: {chair.description}</p>
                            </CardContent>
                          </CardActionArea>
                        </Card>
                      </Link>
                    ))
                  }
                </>
              ) : (
                <CircularProgress />
              )}
            </Box>
          </Container>
        </Paper>
      ) : null}
    </>
  )
}

const ChairSearchPage = () => {
  const [chairRangeMap, setChairRangeMap] = useState<ChairRangeMap | null>(null)

  useEffect(() => {
    fetch('/api/chair/range', { mode: 'cors' })
      .then(async response => await response.json())
      .then(chair => setChairRangeMap(chair as ChairRangeMap))
      .catch(console.error)
  }, [])

  return chairRangeMap ? (
    <ChairSearch chairRangeMap={chairRangeMap} />
  ) : (
    <Loading />
  )
}

export default ChairSearchPage
