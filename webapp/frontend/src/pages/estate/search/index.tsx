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
import { CheckboxForm } from '../../../components/CheckboxForm'

import type { FC } from 'react'
import type { Estate, EstateRangeMap, EstateSearchCondition, EstateSearchResponse } from '@types'

const ESTATE_COUNTS_PER_PAGE = 20

interface EstateItemProps {
  estate: Estate
}

interface EstateSearchProps {
  estateRangeMap: EstateRangeMap
}

const useEstateItemStyles = makeStyles(theme =>
  createStyles({
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

const useEstateSearchStyles = makeStyles(theme =>
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

const FEATURE_LIST = [
  'バストイレ別',
  '駅から徒歩5分',
  'ペット飼育可能',
  'デザイナーズ物件'
]

const EstateItem: FC<EstateItemProps> = ({ estate }) => {
  const classes = useEstateItemStyles()

  return (
    <Link key={estate.id} href={`/estate/detail?id=${estate.id}`}>
      <Card className={classes.card}>
        <CardActionArea className={classes.cardActionArea}>
          <CardMedia image={estate.thumbnail} className={classes.cardMedia} />
          <CardContent className={classes.cardContent}>
            <h2>{estate.name}</h2>
            <p>住所: {estate.address}</p>
            <p>価格: {estate.rent}円</p>
            <p>詳細: {estate.description}</p>
          </CardContent>
        </CardActionArea>
      </Card>
    </Link>
  )
}

const EstateSearch: FC<EstateSearchProps> = ({ estateRangeMap }) => {
  const classes = useEstateSearchStyles()

  const [doorWidthRangeId, setDoorWidthRangeId] = useState('')
  const [doorHeightRangeId, setDoorHeightRangeId] = useState('')
  const [rentRangeId, setRentRangeId] = useState('')
  const [features, setFeatures] = useState<boolean[]>(new Array(FEATURE_LIST.length).fill(false))
  const [estateSearchCondition, setEstateSearchCondition] = useState<EstateSearchCondition | null>(null)
  const [searchResult, setSearchResult] = useState<EstateSearchResponse | null>(null)
  const [page, setPage] = useState<number>(0)

  const onSearch = () => {
    const selectedFeatures = FEATURE_LIST.filter((_, i) => features[i])
    const condition: EstateSearchCondition = {
      doorWidthRangeId,
      doorHeightRangeId,
      rentRangeId,
      features: selectedFeatures.length > 0 ? selectedFeatures.join(',') : '',
      page: 0,
      perPage: ESTATE_COUNTS_PER_PAGE
    }
    setEstateSearchCondition(condition)

    const params = new URLSearchParams()
    for (const [key, value] of Object.entries(condition)) {
      params.append(key, value.toString())
    }
    fetch(`/api/estate/search?${params.toString()}`, { mode: 'cors' })
      .then(async response => await response.json())
      .then(result => {
        setSearchResult(result as EstateSearchResponse)
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
              name='ドアの横幅'
              value={doorWidthRangeId}
              rangeList={estateRangeMap.doorWidth}
              onChange={(_, value) => { setDoorWidthRangeId(value) }}
            />

            <RangeForm
              name='ドアの高さ'
              value={doorHeightRangeId}
              rangeList={estateRangeMap.doorHeight}
              onChange={(event, value) => { setDoorHeightRangeId(event.target.value) }}
            />

            <RangeForm
              name='賃料'
              value={rentRangeId}
              rangeList={estateRangeMap.rent}
              onChange={(_, value) => { setRentRangeId(value) }}
            />

            <CheckboxForm
              name='特徴'
              checkList={features}
              selectList={FEATURE_LIST}
              onChange={(_, checked, key) => {
                setFeatures(
                  features.map((feature, i) => key === i ? checked : feature)
                )
              }}
            />

            <Button
              onClick={onSearch}
              disabled={
                doorWidthRangeId === '' &&
                doorHeightRangeId === '' &&
                rentRangeId === '' &&
                !features.some(feature => feature)
              }
            >
              Search
            </Button>
          </Box>
        </Container>
      </Paper>

      {estateSearchCondition ? (
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
                      if (!estateSearchCondition) return
                      const condition = { ...estateSearchCondition, page: page - 1 }
                      setEstateSearchCondition(condition)
                      setSearchResult(null)
                      const params = new URLSearchParams()
                      for (const [key, value] of Object.entries(condition)) {
                        params.append(key, value.toString())
                      }
                      fetch(`/api/estate/search?${params.toString()}`, { mode: 'cors' })
                        .then(async response => await response.json())
                        .then(result => {
                          setSearchResult(result as EstateSearchResponse)
                        })
                        .catch(console.error)
                    }}
                  />
                  {
                    searchResult.estates.map((estate, i) => (
                      <EstateItem key={i} estate={estate} />
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

const EstateSearchPage = () => {
  const [estateRangeMap, setEstateRangeMap] = useState<EstateRangeMap | null>(null)

  useEffect(() => {
    fetch('/api/estate/range', { mode: 'cors' })
      .then(async response => await response.json())
      .then(estate => setEstateRangeMap(estate as EstateRangeMap))
      .catch(console.error)
  }, [])

  return estateRangeMap ? (
    <EstateSearch estateRangeMap={estateRangeMap} />
  ) : (
    <Loading />
  )
}

export default EstateSearchPage
