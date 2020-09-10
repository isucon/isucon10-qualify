<?php
declare(strict_types=1);

use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Psr\Log\LoggerInterface;
use Fig\Http\Message\StatusCodeInterface;
use League\Csv\Reader;
use Slim\App;

const EXEC_SUCCESS = 127;
const DIR_ROOT = __DIR__ . '/../..';

class Chair
{
    public function getId(): ?int
    {
        return (int)$this->id;
    }

    public function getName(): ?string
    {
        return $this->name;
    }

    public function getDescription(): ?string
    {
        return $this->description;
    }

    public function getThumbnail(): ?string
    {
        return $this->thumbnail;
    }

    public function getPrice(): ?int
    {
        return (int)$this->price;
    }

    public function getHeight(): ?int
    {
        return (int)$this->height;
    }

    public function getWidth(): ?int
    {
        return (int)$this->width;
    }

    public function getDepth(): ?int
    {
        return (int)$this->depth;
    }

    public function getColor(): ?string
    {
        return $this->color;
    }

    public function getFeatures(): ?string
    {
        return $this->features;
    }

    public function getKind(): ?string
    {
        return $this->kind;
    }

    public function getPopularity(): ?int
    {
        return (int)$this->popularity;
    }

    public function getStock(): ?int
    {
        return (int)$this->stock;
    }

    public function toArray()
    {
        return [
            'id' => $this->getId(),
            'name' => $this->getName(),
            'description' => $this->getDescription(),
            'thumbnail' => $this->getThumbnail(),
            'price' => $this->getPrice(),
            'height' => $this->getHeight(),
            'width' => $this->getWidth(),
            'depth' => $this->getDepth(),
            'color' => $this->getColor(),
            'features' => $this->getFeatures(),
            'kind' => $this->getKind(),
            'popularity' => $this->getPopularity(),
            'stock' => $this->getStock(),
        ];
    }
}

class Range
{
    public ?int $id;
    public ?int $min;
    public ?int $max;

    public function __construct(
        int $id = null,
        int $min = null,
        int $max = null
    ) {
        $this->id = $id;
        $this->min = $min;
        $this->max = $max;
    }

    public static function unmarshal(array $json) {
        return new Range(
            $json['id'] ?? null,
            $json['min'] ?? null,
            $json['max'] ?? null,
        );
    }
}

class RangeCondition
{
    public ?string $prefix;
    public ?string $suffix;
    /** @var Range[] */
    public array $ranges;

    public function __construct(
        string $prefix = null,
        string $suffix = null,
        array $ranges = []
    ) {
        $this->prefix = $prefix;
        $this->suffix = $suffix;
        $this->ranges = $ranges;
    }

    public static function unmarshal(array $json): RangeCondition
    {
        return new RangeCondition(
            $json['prefix'] ?? null,
            $json['suffix'] ?? null,
            array_map(Range::class . '::unmarshal', $json['ranges'] ?? [])
        );
    }
}

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

class EstateSearchCondition
{
    public ?RangeCondition $doorWidth;
    public ?RangeCondition $doorHeight;
    public ?RangeCondition $rent;
    public ?ListCondition $feature;

    public function __construct(
        RangeCondition $doorHeight = null,
        RangeCondition $doorWidth = null,
        RangeCondition $rent = null,
        ListCondition $feature = null
    ) {
        $this->doorHeight = $doorHeight;
        $this->doorWidth = $doorWidth;
        $this->rent = $rent;
        $this->feature = $feature;
    }

    public static function unmarshal(array $json): EstateSearchCondition
    {
        return new EstateSearchCondition(
            isset($json['doorHeight']) ? RangeCondition::unmarshal($json['doorHeight']) : null,
            isset($json['doorWidth']) ? RangeCondition::unmarshal($json['doorWidth']) : null,
            isset($json['rent']) ? RangeCondition::unmarshal($json['rent']) : null,
            isset($json['feature']) ? ListCondition::unmarshal($json['feature']) : null,
        );
    }
}

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

function getRange($condition, int $rangeId)
{
    // $rangeId
}

global $chairSearchCondition;
global $estateSearchCondition;

function init() {
    global $chairSearchCondition;
    global $estateSearchCondition;

    if (!$jsonText = file_get_contents(DIR_ROOT . '/fixture/chair_condition.json')) {
        throw new RuntimeException(sprintf('Failed to get load file: %s', '/fixture/chair_condition.json'));
    }
    if (!$json = json_decode($jsonText, true)) {
        throw new RuntimeException(sprintf('Failed to parse json: %s', '../fixture/chair_condition.json'));
    }
    $chairSearchCondition = ChairSearchCondition::unmarshal($json);

    if (!$jsonText = file_get_contents(DIR_ROOT . '/fixture/estate_condition.json')) {
        throw new RuntimeException(sprintf('Failed to get load file: %s', '/fixture/estate_condition.json'));
    }
    if (!$json = json_decode($jsonText, true)) {
        throw new RuntimeException(sprintf('Failed to parse json: %s', '../fixture/estate_condition.json'));
    }
    $estateSearchCondition = EstateSearchCondition::unmarshal($json);
}

init();

return function (App $app) {
    $app->options('/{routes:.*}', function (Request $request, Response $response) {
        // CORS Pre-Flight OPTIONS Request Handler
        return $response;
    });

    $app->post('/initialize', function(Request $request, Response $response): Response {
        $config = $this->get('settings')['database'];

        $paths = [
            '../mysql/db/0_Schema.sql',
            '../mysql/db/1_DummyEstateData.sql',
            '../mysql/db/2_DummyChairData.sql',
        ];

        foreach ($paths as $path) {
            $sqlFile = realpath($path);
            $cmdStr = vsprintf('mysql -h %s -u %s -p%s %s < %s', [
                $config['host'],
                $config['user'],
                $config['pass'],
                $config['dbname'],
                $sqlFile,
            ]);

            system("bash -c $cmdStr", $result);
            if ($result !== EXEC_SUCCESS) {
                $this->get(LoggerInterface::class)->error('Initialize script error');
                return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
            }
        }

        $response->getBody()->write(json_encode([
            'language' => 'php',
        ]));

        return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(200);
    });

    $app->get('/api/chair/search', function(Request $request, Response $response) {
        $conditions = [];
        $params = [];

        if ($priceRangeId = $request->getQueryParams()['priceRangeId'] ?? null) {
            if (!is_numeric($priceRangeId)) {
                $this->get(LoggerInterface::class)->error(sprintf('priceRangeID invalid, %s', $priceRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            $chairPrice = getRange($condition, (int)$rangeId);
        }

        return $response;
    });

    $app->get("/api/chair/{id}", function(Request $request, Response $response, array $args) {
        $id = $args['id'] ?? null;
        if (empty($id) || !is_numeric($id)) {
            $this->get(LoggerInterface::class)->error(sprintf('Request parameter \"id\" parse error : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
        }

        $query = 'SELECT * FROM chair WHERE id = :id';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->execute([':id' => $id]);
        $chair = $stmt->fetchObject(Chair::class);

        if (!$chair) {
            $this->get(LoggerInterface::class)->error(sprintf('requested id\'s chair not found : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
        } elseif (!$chair instanceof Chair) {
            $this->get(LoggerInterface::class)->error(sprintf('Failed to get the chair from id : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
        } elseif ($chair->getStock() <= 0) {
            $this->get(LoggerInterface::class)->error(sprintf('requested id\'s chair is sold out : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
        }

        $response->getBody()->write(json_encode($chair->toArray()));

        return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(200);
    });

    $app->post('/api/chair', function (Request $request, Response $response) {
        if (!$file = $request->getUploadedFiles()['chairs'] ?? null) {
            $this->get(LoggerInterface::class)->error('failed to get form file');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        } elseif (!$file instanceof Slim\Psr7\UploadedFile || $file->getError() !== UPLOAD_ERR_OK) {
            $this->get(LoggerInterface::class)->error('failed to get form file');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        if (!$records = Reader::createFromPath($file->getFilePath())) {
            $this->get(LoggerInterface::class)->error('failed to read csv');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        $pdo = $this->get(PDO::class);

        if (!$pdo->beginTransaction()) {
            $this->get(LoggerInterface::class)->error('failed to begin tx');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        try {
            foreach ($records as $record) {
                $query = 'INSERT INTO chair VALUES(:id, :name, :description, :thumbnail, :price, :height, :width, :depth, :color, :features, :kind, :popularity, :stock)';
                $stmt = $pdo->prepare($query);
                $stmt->execute([
                    ':id' => (int)trim($record[0] ?? null),
                    ':name' => (string)trim($record[1] ?? null),
                    ':description' => (string)trim($record[2] ?? null),
                    ':thumbnail' => (string)trim($record[3] ?? null),
                    ':price' => (int)trim($record[4] ?? null),
                    ':height' => (int)trim($record[5] ?? null),
                    ':width' => (int)trim($record[6]) ?? null,
                    ':depth' => (int)trim($record[7] ?? null),
                    ':color' => (string)trim($record[8] ?? null),
                    ':features' => (string)trim($record[9] ?? null),
                    ':kind' => (string)trim($record[10] ?? null),
                    ':popularity' => (int)trim($record[11] ?? null),
                    ':stock' => (int)trim($record[12] ?? null),
                ]);
            }
        } catch (PDOException $e) {
            $pdo->rollback();
            $this->get(LoggerInterface::class)->error(sprintf('failed to insert chair: %s', $e->getMessage()));
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        if (!$pdo->commit()) {
            $this->get(LoggerInterface::class)->error('failed to commit tx');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
    });

    $app->get('/', function (Request $request, Response $response) {
        $response->getBody()->write('Hello world!');
        return $response;
    });
};
