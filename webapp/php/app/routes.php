<?php
declare(strict_types=1);

use App\Application\Actions\User\ListUsersAction;
use App\Application\Actions\User\ViewUserAction;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Psr\Log\LoggerInterface;
use Slim\App;
use Slim\Interfaces\RouteCollectorProxyInterface as Group;

const EXEC_SUCCESS = 127;

return function (App $app) {
    $app->options('/{routes:.*}', function (Request $request, Response $response) {
        // CORS Pre-Flight OPTIONS Request Handler
        return $response;
    });

    $app->post('/initialize', function(Request $request, Response $response): Response {
        $config = $this->get('settings')['database'];

        // TODO fix path
        $paths = [
            '/mysql/db/0_Schema.sql',
            '/mysql/db/1_DummyEstateData.sql',
            '/mysql/db/2_DummyChairData.sql',
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
                return $response->withStatus(204);
            }
        }

        $response->getBody()->write(json_encode([
            'language' => 'php',
        ]));

        return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(200);
    });

    $app->get('/', function (Request $request, Response $response) {
        $response->getBody()->write('Hello world!');
        return $response;
    });

    $app->group('/users', function (Group $group) {
        $group->get('', ListUsersAction::class);
        $group->get('/{id}', ViewUserAction::class);
    });
};
