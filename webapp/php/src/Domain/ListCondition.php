<?php

namespace App\Domain;

class ListCondition
{
    /** @var string[] */
    public array $list;

    public function __construct(array $list)
    {
        $this->list = $list;
    }

    public static function unmarshal(array $list): ListCondition
    {
        return new ListCondition($list);
    }
}