package main

import (
	"github.com/webability-go/xamboo/assets"
	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xmodules/user/bridge"
)

func Run(ctx *assets.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	ok := bridge.Setup(ctx, "", bridge.USER)
	if !ok {
		return ""
	}

	params := &xcore.XDataset{
		"#": language,
	}

	return template.Execute(params)
}

/*
<?php

class derechos extends \common\WAApplication
{
  private $usuarioEntity;

  public function __construct($template, $language)
  {
    parent::__construct($template, $language);

    $this->usuarioEntity = \entities\usuarioEntity::getInstance();
  }

  public function code($engine, $params)
  {
    $code = $this->template->resolve();
    return $code;
  }

  public function contenedoresdata()
  {
    $init = microtime(true);
    $data = null;
    if (isset($_POST['data']))
      $data = json_decode($_POST['data']);

    if (!$data)
    {
      $data = array(array(0, 49));
    }

    $t1 = microtime(true);
    $total = $this->usuarioEntity->countAdminDerecho(null);
    $t2 = microtime(true);
    $rec = array('total' => $total, 'row' => array());
    $rec['subtimetotal'] = $t2 - $t1;

    // condicionar por orden y lista
    foreach($data as $set)
    {
      $min = $set[0];
      $max = $set[1];

      $t1 = microtime(true);
      $entradas = $this->usuarioEntity->readAdminDerecho(null, new \dominion\DB_OrderBy('nombre', \dominion\DB_OrderBy::ASC), $max-$min+1, $min);
      $t2 = microtime(true);
      $rec['subtime'] = $t2 - $t1;
      $num = $set[0];

      if ($entradas)
      {
        $turn = true;
        foreach($entradas as $entrada)
        {
          // Busca perfiles con este derechos
          $perfiles = $this->searchperfiles($entrada->clave);
          $usuarios = $this->searchusuarios($entrada->clave);

          $botones = <<<EOF
  <input type="button" value="Modificar" onclick="derechos_editar('{$entrada->clave}');" />
  <input type="button" value="Borrar" onclick="derechos_borrar('{$entrada->clave}');" />
EOF;
          $rec['row'][$num++] = array('clave' => $entrada->clave,
          'nombre' => $entrada->nombre,
          'perfiles' => $perfiles,
          'usuarios' => $usuarios,
          'comandos' => $botones
          );
          $turn = !$turn;
        }
      }
    }
    // flag "fullload" is not set since we load only a part of the dynamic data
    $end = microtime(true);
    $rec['time'] = $end - $init;

    return $rec;
  }

  private function searchperfiles($key)
  {
    $resultado = array();
    $recs = $this->usuarioEntity->selectAdminperfilderecho(array(new \dominion\DB_Condition('derecho', '=', $key)), new \dominion\DB_OrderBy('perfil', \dominion\DB_OrderBy::ASC));
    if ($recs)
    {
      foreach($recs as $rec)
      {
        $perfil = $this->usuarioEntity->selectAdminperfil($rec->perfil);
        $resultado[] = '[<b>' . $perfil->clave.'</b>]';
      }
    }
    if (!$resultado)
      return '--';

    return implode(', ', $resultado);
  }

  private function searchusuarios($key)
  {
    $resultado = array();
    $usuarios = array();

    // usuarios directos con este derecho SI/NO
    $recs = $this->usuarioEntity->selectAdminderechousuario(array(new \dominion\DB_Condition('derecho', '=', $key)), new \dominion\DB_OrderBy('usuario', \dominion\DB_OrderBy::ASC));
    if ($recs)
    {
      foreach($recs as $rec)
      {
        $usuario = $this->usuarioEntity->selectUsuario($rec->usuario);
        if ($usuario->estatus == 'X')
          continue;
        if ($rec->estatus == 1)
          $color = 'green';
        elseif ($rec->estatus == 2)
          $color = 'red';
        else
          $color = 'black';

        $usuarios[$usuario->nombre] = $color;
      }
    }

    // usuarios indirectos a traves del perfil
    $recs = $this->usuarioEntity->selectAdminperfilderecho(array(new \dominion\DB_Condition('derecho', '=', $key)), new \dominion\DB_OrderBy('perfil', \dominion\DB_OrderBy::ASC));
    if ($recs)
    {
      foreach($recs as $rec)
      {
        $subrecs = $this->usuarioEntity->selectAdminPerfilUsuario(array(new \dominion\DB_Condition('perfil', '=', $rec->perfil)), new \dominion\DB_OrderBy('usuario', \dominion\DB_OrderBy::ASC));
        if ($subrecs)
        {
          foreach($subrecs as $subrec)
          {
            $usuario = $this->usuarioEntity->selectUsuario($subrec->usuario);
            if ($usuario->estatus == 'X')
              continue;
            if ($rec->estatus == 1)
              $color = 'green';
            elseif ($rec->estatus == 2)
              $color = 'red';
            else
              $color = 'black';
            $usuarios[$usuario->nombre] = $color;
          }
        }
      }
    }

    if (!$usuarios)
      return '--';

    foreach($usuarios as $usuario => $color)
      $resultado[] = '<span style="color: '.$color.';">[' . $usuario.']</span>';

    return implode(', ', $resultado);
  }

  public function synchro()
  {
    // scan menu
    include_once($this->base->config->BASEDIR . 'application/pagesadmin/main/main.lib');
    $main = new main(null, null);
    $menu = $main->getMenuData();

    $rights = array();
    $num = 0; // num '0' no existira: la primera opcion es abierta sin derecho e incrementa a 1
    foreach($menu as $group)
    {
      if ($group['right'])
        $rights[$group['right']] = ($num<10?'0':'') . $num . '.00 ' . $group['name'];
      $subnum = 1;
      if (isset($group['pages']))
      {
        foreach($group['pages'] as $option)
        {
          if ($option['right'])
            $rights[$option['right']] = ($num<10?'0':'') . $num . '.' . ($subnum<10?'0':'') . $subnum . ' ' . $option['name'];
          $subnum++;
        }
      }
      $num++;
    }

    foreach($rights as $right => $name)
    {
      $exist = $this->usuarioEntity->selectAdminDerecho($right);
      if ($exist)
        $this->usuarioEntity->updateAdminDerecho($right, array('nombre' => $name));
      else
        $this->usuarioEntity->insertAdminDerecho(array('clave' => $right, 'nombre' => $name));
    }

    return array('estatus' => 'OK');
  }

}

?>
*/
