<?php

namespace App\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;
use App\Repository\CountRepository;

class HomeController extends AbstractController
{
    #[Route('/')]
    public function index(): Response 
    {
        return $this->render('pages/index.html.twig');
    }
}