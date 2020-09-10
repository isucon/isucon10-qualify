<?php
declare(strict_types=1);

use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Psr\Log\LoggerInterface;
use Fig\Http\Message\StatusCodeInterface;
use Slim\App;

const EXEC_SUCCESS = 127;

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

    $app->get('/', function (Request $request, Response $response) {
        $response->getBody()->write('Hello world!');
        return $response;
    });
};
