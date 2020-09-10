<?php
declare(strict_types=1);

use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Psr\Log\LoggerInterface;
use Fig\Http\Message\StatusCodeInterface;
use League\Csv\Reader;
use Slim\App;
use App\Domain\Chair;

const EXEC_SUCCESS = 127;

function getRange($condition, int $rangeId)
{
    // $rangeId
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
