<?php
declare(strict_types=1);

use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Psr\Log\LoggerInterface;
use Fig\Http\Message\StatusCodeInterface;
use League\Csv\Reader;
use Slim\App;
use App\Domain\Chair;
use App\Domain\ChairSearchCondition;
use App\Domain\Range;
use App\Domain\RangeCondition;

const EXEC_SUCCESS = 127;
const NUM_LIMIT = 20;

function getRange(RangeCondition $condition, int $rangeId): ?Range
{
    if ($rangeId < 0) {
        return null;
    }
    if (count($condition->ranges) <= $rangeId) {
        return null;
    }

    return $condition->ranges[$rangeId] ?? null;
}

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

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/search', function(Request $request, Response $response) {
        $conditions = [];
        $params = [];

        /** @var ChairSearchCondition */
        $chairSearchCondition = $this->get(ChairSearchCondition::class);

        if ($priceRangeId = $request->getQueryParams()['priceRangeId'] ?? null) {
            $chairPrice = getRange($chairSearchCondition->price, (int)$priceRangeId);
            if (!$chairPrice) {
                $this->get(LoggerInterface::class)->info(sprintf('priceRangeId invalid, %s', $priceRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairPrice->min != -1) {
                $conditions[] = 'price >= :minPrice';
                $params[':minPrice'] = [$chairPrice->min, PDO::PARAM_INT];
            }
            if ($chairPrice->max != -1) {
                $conditions[] = 'price < :maxPrice';
                $params[':maxPrice'] = [$chairPrice->max, PDO::PARAM_INT];
            }
        }
        if ($heightRangeId = $request->getQueryParams()['heightRangeId'] ?? null) {
            $chairHeight = getRange($chairSearchCondition->height, $heightRangeId);
            if (!$chairHeight) {
                $this->get(LoggerInterface::class)->info(sprintf('heightRangeId invalid, %s', $heightRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairHeight->min != -1) {
                $conditions[] = 'height >= :minHeight';
                $params[':minHeight'] = [$chairHeight->min, PDO::PARAM_INT];
            }
            if ($chairHeight->max != -1) {
                $conditions[] = 'height < :maxHeight';
                $params[':maxHeight'] = [$chairHeight->max, PDO::PARAM_INT];
            }
        }
        if ($widthRangeId = $request->getQueryParams()['widthRangeId'] ?? null) {
            $chairWidth = getRange($chairSearchCondition->width, $widthRangeId);
            if (!$chairWidth) {
                $this->get(LoggerInterface::class)->info(sprintf('widthRangeId invalid, %s', $heightRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairWidth->min != -1) {
                $conditions[] = 'width >= :minWidth';
                $params[':minWidth'] = [$chairWidth->min, PDO::PARAM_INT];
            }
            if ($chairWidth->max != -1) {
                $conditions[] = 'width < :maxWidth';
                $params[':maxWidth'] = [$chairWidth->max, PDO::PARAM_INT];
            }
        }
        if ($depthRangeId = $request->getQueryParams()['depthRangeId'] ?? null) {
            $chairDepth = getRange($chairSearchCondition->depth, $depthRangeId);
            if (!$chairDepth) {
                $this->get(LoggerInterface::class)->info(sprintf('depthRangeId invalid, %s', $heightRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairDepth->min != -1) {
                $conditions[] = 'depth >= :minDepth';
                $params[':minDepth'] = [$chairDepth->min, PDO::PARAM_INT];
            }
            if ($chairDepth->max != -1) {
                $conditions[] = 'depth < :maxDepth';
                $params[':maxDepth'] = [$chairDepth->max, PDO::PARAM_INT];
            }
        }
        if ($kind = $request->getQueryParams()['kind'] ?? null) {
            $conditions[] = 'kind = :kind';
            $params[':kind'] = [$kind, PDO::PARAM_STR];
        }
        if ($color = $request->getQueryParams()['color'] ?? null) {
            $conditions[] = 'color = :color';
            $params[':color'] = [$color, PDO::PARAM_STR];
        }
        if ($features = $request->getQueryParams()['features'] ?? null) {
            foreach (explode(',', $features) as $key => $feature) {
                $name = sprintf(':feature_%s', $key);
                $conditions[] = sprintf("features LIKE CONCAT('%%', %s, '%%')", $name);
                $params[$name] = [$feature, PDO::PARAM_STR];
            }
        }

        if (count($conditions) === 0) {
            $this->get(LoggerInterface::class)->info('Search condition not found');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $conditions[] = 'stock > 0';

        if (is_null($page = $request->getQueryParams()['page'] ?? null)) {
            $this->get(LoggerInterface::class)->info(sprintf('Invalid format page parameter: %s', $page));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }
        if (is_null($perPage = $request->getQueryParams()['perPage'] ?? null)) {
            $this->get(LoggerInterface::class)->info(sprintf('Invalid format perPage parameter: %s', $perPage));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $searchQuery = 'SELECT * FROM chair WHERE ';
        $countQuery = 'SELECT COUNT(*) FROM chair WHERE ';
        $searchCondition = implode(' AND ', $conditions);
        $limitOffset = ' ORDER BY popularity DESC, id ASC LIMIT :limit OFFSET :offset';

        $stmt = $this->get(PDO::class)->prepare($countQuery . $searchCondition);
        foreach ($params as $key => $bind) {
            list($value, $type) = $bind;
            $stmt->bindValue($key, $value, $type);
        }
        $stmt->execute();
        $count = (int)$stmt->fetchColumn();

        $params[':limit'] = [(int)$perPage, PDO::PARAM_INT];
        $params[':offset'] = [(int)$page*$perPage, PDO::PARAM_INT];

        $stmt = $this->get(PDO::class)->prepare($searchQuery . $searchCondition . $limitOffset);
        foreach ($params as $key => $bind) {
            list($value, $type) = $bind;
            $stmt->bindValue($key, $value, $type);
        }
        $stmt->execute();
        $chairs = $stmt->fetchAll(PDO::FETCH_CLASS, Chair::class);

        if (count($chairs) === 0) {
            $response->getBody()->write(json_encode([
                'count' => $count,
                'chairs' => [],
            ]));
            return $response->withHeader('Content-Type', 'application/json');
        }

        $response->getBody()->write(json_encode([
            'count' => $count,
            'chairs' => array_map(
                function(Chair $chair) {
                    return $chair->toArray();
                },
                $chairs
            ),
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/low_priced', function(Request $request, Response $response) {
        $query = 'SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT :limit';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->bindValue(':limit', NUM_LIMIT, PDO::PARAM_INT);
        $stmt->execute();
        $chairs = $stmt->fetchAll(PDO::FETCH_CLASS, Chair::class);

        if (count($chairs) === 0) {
            $this->get(LoggerInterface::class)->error('getLowPricedChair not found');
            $response->getBody()->write(json_encode([
                'chairs' => []
            ]));
            return $response->withHeader('Content-Type', 'application/json');
        }

        $response->getBody()->write(json_encode([
            'chairs' => array_map(
                function(Chair $chair) {
                    return $chair->toArray();
                },
                $chairs
            )
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/search/condition', function(Request $request, Response $response) {
        $chairSearchCondition = $this->get(ChairSearchCondition::class);
        $response->getBody()->write(json_encode($chairSearchCondition));
        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->post('/api/chair/buy/{id}', function(Request $request, Response $response, array $args) {
        if (!$id = $args['id'] ?? null) {
            $this->get(LoggerInterface::class)->info('post request document failed');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $pdo = $this->get(PDO::class);

        try {
            $pdo->beginTransaction();
            $stmt = $pdo->prepare('SELECT * FROM chair WHERE id = :id AND stock > 0 FOR UPDATE');
            $stmt->bindValue(':id', $id, PDO::PARAM_INT);
            $stmt->execute();
            $chair = $stmt->fetchObject(Chair::class);

            if (!$chair) {
                $pdo->rollback();
                $this->get(LoggerInterface::class)->info(sprintf('buyChair chair id "%s" not found', $id));
                return $response->withStatus(StatusCodeInterface::STATUS_NOT_FOUND);
            }

            $stmt = $pdo->prepare('UPDATE chair SET stock = stock - 1 WHERE id = :id');
            $stmt->bindValue(':id', $id, PDO::PARAM_INT);
            if (!$stmt->execute()) {
                $pdo->rollback();
                $this->get(LoggerInterface::class)->error('chair stock update failed');
                return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
            }

            if (!$pdo->commit()) {
                $this->get(LoggerInterface::class)->error('transaction commit error');
                return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
            }
        } catch (PDOException $e) {
            $pdo->rollBack();
            $this->get(LoggerInterface::class)->error(sprintf('DB Execution Error: on getting a chair by id : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/{id}', function(Request $request, Response $response, array $args) {
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

        return $response->withHeader('Content-Type', 'application/json');
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
            if (!$pdo->commit()) {
                $this->get(LoggerInterface::class)->error('failed to commit tx');
                return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
            }
        } catch (PDOException $e) {
            $pdo->rollBack();
            $this->get(LoggerInterface::class)->error(sprintf('failed to insert chair: %s', $e->getMessage()));
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
    });

    $app->get('/', function (Request $request, Response $response) {
        $response->getBody()->write('Hello world!');
        return $response;
    });
};
