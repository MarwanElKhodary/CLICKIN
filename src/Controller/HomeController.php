<?php

namespace App\Controller;

use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;
use App\Repository\CountRepository;

class HomeController 
{
    #[Route('/')]
    public function index(CountRepository $respository): Response 
    {
        $count = $respository->findAll();

        dd($count);

        return new Response("Hello from a controller!");
    }
}