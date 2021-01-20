package main

import (
	"github.com/webability-go/xcore/v2"

	// We include the master/app bridge, where the xmodules are supposed to live
	// You may change this to you own APP and bridge upon needs
	"github.com/webability-go/xamboo/applications/xmodules/app/bridge"
	"github.com/webability-go/xamboo/assets"
)

func Run(ctx *assets.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	bridge.Setup(ctx)

	//	bridge.EntityLog_LogStat(ctx)
	params := &xcore.XDataset{
		"#": language,
	}

	return template.Execute(params)
}

func Userslist(ctx *assets.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	bridge.Setup(ctx)

	init := time.Now()
	data := [][]int{}
	fdata := ctx.Request.Form.Get("data")
	if fdata != "" {
		json.Unmarshal([]byte(fdata), &data)
	}
	if len(data) == 0 {
		data = append(data, []int{0, 49})
	}

	query := ""
  filtroestatus := bridge.GetUserParamInt("filtroestatus")
	if filtroestatus == 1	{
		query = array(new \dominion\DB_Condition('estatus', 'in', "('A', 'S')"));
	}

	t1 := time.Now();
	totalusers := bridge.GetUsersListCount("xmodules", query)
	t2 := time.Now();

	result := map[string]interface{}{
		"total": totalusers,
		"row": map[int]interface{},
		"subtimetotal": t2-t1,
	}

	// condicionar por orden y lista
	for _, set := range data {
		min = set.min
		max = set.max

		$t1 = microtime(true);
		$entradas = $this->usuarioEntity->readUsuario($query, new \dominion\DB_OrderBy('clave', \dominion\DB_OrderBy::ASC), $max-$min+1, $min);
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
				$derechos = $this->searchderechos($entrada->clave);

				$botones = <<<EOF
<input type="button" value="Modificar" onclick="usuarios_editar('{$entrada->clave}');" />
<input type="button" value="Borrar" onclick="usuarios_borrar('{$entrada->clave}');" />
EOF;
				$rec['row'][$num++] = array('clave' => $entrada->clave,
				'estatus' => $entrada->estatus,
				'nombre' => $entrada->nombre . ' [' . $entrada->chef . ']<br />' . $entrada->correo,
				'perfiles' => $perfiles,
				'derechos' => $derechos,
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







	return map[string]interface{}{
		"row": users,
	}
}

/*
<?php

class usuarios extends \common\WAApplication
{
  private $usuarioEntity;

  private $FILTROESTATUS;

  public function __construct($template, $language)
  {
    parent::__construct($template, $language);

    $this->usuarioEntity = \entities\usuarioEntity::getInstance();

    $this->FILTROESTATUS = $this->usuarioEntity->getParam('USUARIOS_ESTATUS');
    if (!$this->FILTROESTATUS)
    {
      $this->FILTROESTATUS = 1;
      $this->usuarioEntity->setParam('USUARIOS_ESTATUS', $this->FILTROESTATUS);
    }
  }

  public function code($engine, $params)
  {
    $check = '';
    if ($this->FILTROESTATUS == 1)
    {
      $check = ' checked="checked"';
    }

    $this->template->metaElements(
      array(
        'checkfiltro' => $check,
      )
    );
    $code = $this->template->resolve();
    return $code;
  }

  public function filtro($engine, $params)
  {
    if ($this->FILTROESTATUS == 1)
      $this->FILTROESTATUS = 2;
    else
      $this->FILTROESTATUS = 1;
    $this->usuarioEntity->setParam('USUARIOS_ESTATUS', $this->FILTROESTATUS );

    return array('estatus' => 'OK', 'filtro' => $this->FILTROESTATUS);
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

    $query = null;
    if ($this->FILTROESTATUS == 1)
    {
      $query = array(new \dominion\DB_Condition('estatus', 'in', "('A', 'S')"));
    }

    $t1 = microtime(true);
    $total = $this->usuarioEntity->countUsuario($query);
    $t2 = microtime(true);
    $rec = array('total' => $total, 'row' => array());
    $rec['subtimetotal'] = $t2 - $t1;

    // condicionar por orden y lista
    foreach($data as $set)
    {
      $min = $set[0];
      $max = $set[1];

      $t1 = microtime(true);
      $entradas = $this->usuarioEntity->readUsuario($query, new \dominion\DB_OrderBy('clave', \dominion\DB_OrderBy::ASC), $max-$min+1, $min);
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
          $derechos = $this->searchderechos($entrada->clave);

          $botones = <<<EOF
  <input type="button" value="Modificar" onclick="usuarios_editar('{$entrada->clave}');" />
  <input type="button" value="Borrar" onclick="usuarios_borrar('{$entrada->clave}');" />
EOF;
          $rec['row'][$num++] = array('clave' => $entrada->clave,
          'estatus' => $entrada->estatus,
          'nombre' => $entrada->nombre . ' [' . $entrada->chef . ']<br />' . $entrada->correo,
          'perfiles' => $perfiles,
          'derechos' => $derechos,
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
    $recs = $this->usuarioEntity->selectAdminperfilusuario(array(new \dominion\DB_Condition('usuario', '=', $key)), new \dominion\DB_OrderBy('perfil', \dominion\DB_OrderBy::ASC));
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

  private function searchderechos($key)
  {
    $resultado = array();
    $derechos = array();

    // derechos directos con este derecho SI/NO
    $recs = $this->usuarioEntity->selectAdminderechousuario(array(new \dominion\DB_Condition('usuario', '=', $key)), new \dominion\DB_OrderBy('derecho', \dominion\DB_OrderBy::ASC));
    if ($recs)
    {
      foreach($recs as $rec)
      {
        if ($rec->estatus == 1)
          $color = 'green';
        elseif ($rec->estatus == 2)
          $color = 'red';
        else
          $color = 'black';

        $derechos[$rec->derecho] = $color;
      }
    }

    // derechos indirectos a traves del perfil
    $recs = $this->usuarioEntity->selectAdminperfilusuario(array(new \dominion\DB_Condition('usuario', '=', $key)), new \dominion\DB_OrderBy('perfil', \dominion\DB_OrderBy::ASC));
    if ($recs)
    {
      foreach($recs as $rec)
      {
        $subrecs = $this->usuarioEntity->selectAdminPerfilderecho(array(new \dominion\DB_Condition('perfil', '=', $rec->perfil)), new \dominion\DB_OrderBy('derecho', \dominion\DB_OrderBy::ASC));
        if ($subrecs)
        {
          foreach($subrecs as $subrec)
          {
            if ($subrec->estatus == 1)
              $color = 'green';
            elseif ($subrec->estatus == 2)
              $color = 'red';
            else
              $color = 'black';
            $derechos[$subrec->derecho] = $color;
          }
        }
      }
    }

    if (!$derechos)
      return '--';

    foreach($derechos as $derecho => $color)
      $resultado[] = '<span style="color: '.$color.';">[' . $derecho.']</span>';

    return implode(', ', $resultado);
  }
}

?>
*/
