<?php

namespace App\Domain;

class ChairSearchCondition
{
    public ?RangeCondition $width;
    public ?RangeCondition $height;
    public ?RangeCondition $depth;
    public ?RangeCondition $price;
    public ?RangeCondition $color;
    public ?RangeCondition $feature;
    public ?RangeCondition $kind;

    public function __construct(
        RangeCondition $width = null,
        RangeCondition $height = null,
        RangeCondition $depth = null,
        RangeCondition $price = null,
        RangeCondition $color = null,
        RangeCondition $feature = null,
        RangeCondition $kind = null
    ) {
        $this->width = $width;
        $this->height = $height;
        $this->depth = $depth;
        $this->price = $price;
        $this->color = $color;
        $this->feature = $feature;
        $this->kind = $kind;
    }

    public static function unmarshal(array $json): ChairSearchCondition
    {
        return new ChairSearchCondition(
            isset($json['width']) ? RangeCondition::unmarshal($json['width']) : null,
            isset($json['height']) ? RangeCondition::unmarshal($json['height']) : null,
            isset($json['depth']) ? RangeCondition::unmarshal($json['depth']) : null,
            isset($json['price']) ? RangeCondition::unmarshal($json['price']) : null,
            isset($json['color']) ? RangeCondition::unmarshal($json['color']) : null,
            isset($json['features']) ? RangeCondition::unmarshal($json['features']) : null,
            isset($json['kind']) ? RangeCondition::unmarshal($json['kind']) : null,
        );
    }
}