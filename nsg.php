<?php
// Converts XML to JSON. Nightmare.

date_default_timezone_set('Asia/Singapore');
header('Content-Type: application/json');

$creds = parse_ini_file(".creds.ini");
$url = "http://api.nea.gov.sg/api/WebAPI/?dataset=psi_update&keyref=" . $creds["key"];

$ch = curl_init();
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_URL,$url);
curl_setopt($ch, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V4 );
$result=curl_exec($ch);
$info = curl_getinfo($ch);
$errinfo = curl_error($ch);
curl_close($ch);

$dir = 'log/' . date("Y-m-d");
$fn = $dir . '/' . time() . ".txt";
@mkdir($dir, 0777, true);

file_put_contents($fn, "[Curl info]\n" . print_r($info, true) . "\n[Curl Error]\n" . print_r($errinfo, true) . "\n[XML]\n" . $result);

$xml = simplexml_load_string($result);

$array = array();

foreach ($xml->item->region as $region) {
	if ($region->id == "rNO") {
		//print_r($region->record);
		$attr = $region->record->attributes();
		foreach($attr as $key=>$val){
			$array[(string)$key] = date("c", strtotime ((string)$val));
		}
		foreach ($region->record->reading as $reading) {
			$attr = $reading->attributes();
			$t = (string) $attr->type;
			$v = (string) $attr->value;
			$array[$t] = $v;
		}
		break;
	}
}

echo json_encode($array, JSON_PRETTY_PRINT);

?>
