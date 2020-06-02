import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Container, Paper, Box, Button } from '@material-ui/core'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import EventSeatIcon from '@material-ui/icons/EventSeat'
import HouseIcon from '@material-ui/icons/House'
import TouchAppIcon from '@material-ui/icons/TouchApp'

import { EstateCard } from '../components/EstateCard'
import { ChairCard } from '../components/ChairCard'

import type { Estate, Chair } from '@types'

const useStyles = makeStyles(theme =>
  createStyles({
    paper: {
      margin: theme.spacing(2),
      padding: theme.spacing(4)
    },
    link: {
      margin: theme.spacing(2)
    },
    cards: {
      overflowX: 'scroll',
      display: 'flex'
    }
  })
)

const TopPage = () => {
  const classes = useStyles()

  const [recommendedEstates, setRecommendedEstates] = useState<Estate[] | null>(null)
  const [recommendedChairs, setRecommendedChairs] = useState<Chair[] | null>(null)

  useEffect(() => {
    fetch(`${process.env.API_SERVER_NAME ?? ''}/api/recommended_estate`, { mode: 'cors' })
      .then(async response => await response.json())
      .then(json => setRecommendedEstates(json.estates as Estate[]))
      .catch(console.error)

    fetch(`${process.env.API_SERVER_NAME ?? ''}/api/recommended_chair`, { mode: 'cors' })
      .then(async response => await response.json())
      .then(json => setRecommendedChairs(json.chairs as Chair[]))
      .catch(console.error)
  }, [])

  return (
    <Container maxWidth='md'>
      <h1> isuumo </h1>
      <Paper className={classes.paper}>
        <h2> イス・物件を探す </h2>
        <Link href='/chair/search'>
          <Button variant='contained' color='primary' className={classes.link}>
            <EventSeatIcon /> イス検索
          </Button>
        </Link>
        <Link href='/estate/search'>
          <Button variant='contained' color='primary' className={classes.link}>
            <HouseIcon /> 物件検索
          </Button>
        </Link>
        <Link href='/estate/nazotte'>
          <Button variant='contained' color='primary' className={classes.link}>
            <TouchAppIcon /> なぞって検索
          </Button>
        </Link>
      </Paper>
      {recommendedEstates && (
        <Paper className={classes.paper}>
          <h2> オススメの物件 </h2>
          <Box className={classes.cards}>
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
            {recommendedEstates.map(estate => <EstateCard key={estate.id} estate={estate} />)}
          </Box>
        </Paper>
      )}
      {recommendedChairs && (
        <Paper className={classes.paper}>
          <h2> オススメのイス </h2>
          <Box className={classes.cards}>
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
            {recommendedChairs.map(chair => <ChairCard key={chair.id} chair={chair} />)}
          </Box>
        </Paper>
      )}
    </Container>
  )
}

export default TopPage
