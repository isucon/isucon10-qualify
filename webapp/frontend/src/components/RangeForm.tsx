import {
  FormControl,
  FormLabel,
  FormControlLabel,
  Radio,
  RadioGroup
} from '@material-ui/core'

import type { FC, ChangeEvent } from 'react'
import type { RangeList } from '@types'

interface Props {
  name: string
  value: string
  rangeList: RangeList
  onChange: (event: ChangeEvent<HTMLInputElement>, value: string) => void
}

export const RangeForm: FC<Props> = ({ name, value, rangeList, onChange }) => (
  <FormControl component='fieldset'>
    <FormLabel component='legend'>{name}</FormLabel>
    <RadioGroup
      aria-label={name}
      name={name}
      value={value}
      onChange={onChange}
      row
    >
      {
        rangeList.ranges.map(({ id, min, max }) => {
          const minLabel = min !== -1 ? `${rangeList.prefix}${min}${rangeList.suffix} ` : ''
          const maxLabel = max !== -1 ? ` ${rangeList.prefix}${max}${rangeList.suffix}` : ''
          return <FormControlLabel key={id} value={id.toString()} control={<Radio />} label={`${minLabel}〜${maxLabel}`} />
        })
      }
      <FormControlLabel value='' control={<Radio />} label='指定なし' />
    </RadioGroup>
  </FormControl>
)
