package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	_ "modernc.org/sqlite"
)

type ServingMeasure struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Food struct {
	Codigo                      string           `json:"codigo"`
	Nome                        string           `json:"nome"`
	Slug                        string           `json:"slug"`
	OutrasMedidas               []ServingMeasure `json:"outrasMedidas"`
	Grupo                       string           `json:"grupo"`
	Source                      string           `json:"source"` // "taco" or "tbca"
	EnergiaKJ                   *float64         `json:"energiaKJ"`
	EnergiaKcal                 *float64         `json:"energiaKcal"`
	Umidade                     *float64         `json:"umidade"`
	CarboidratoTotal            *float64         `json:"carboidratoTotal"`
	Proteina                    *float64         `json:"proteina"`
	Lipidios                    *float64         `json:"lipidios"`
	FibraAlimentar              *float64         `json:"fibraAlimentar"`
	Cinzas                      *float64         `json:"cinzas"`
	Colesterol                  *float64         `json:"colesterol"`
	AcidosGraxosSaturados       *float64         `json:"acidosGraxosSaturados"`
	AcidosGraxosMonoinsaturados *float64         `json:"acidosGraxosMonoinsaturados"`
	AcidosGraxosPoliinsaturados *float64         `json:"acidosGraxosPoliinsaturados"`
	Calcio                      *float64         `json:"calcio"`
	Ferro                       *float64         `json:"ferro"`
	Sodio                       *float64         `json:"sodio"`
	Magnesio                    *float64         `json:"magnesio"`
	Fosforo                     *float64         `json:"fosforo"`
	Potassio                    *float64         `json:"potassio"`
	Manganes                    *float64         `json:"manganes"`
	Zinco                       *float64         `json:"zinco"`
	Cobre                       *float64         `json:"cobre"`
	VitaminaARE                 *float64         `json:"vitaminaARE"`
	VitaminaARAE                *float64         `json:"vitaminaARAE"`
	Tiamina                     *float64         `json:"tiamina"`
	Riboflavina                 *float64         `json:"riboflavina"`
	Niacina                     *float64         `json:"niacina"`
	VitaminaC                   *float64         `json:"vitaminaC"`
	// TACO specific
	Retinol    *float64 `json:"retinol,omitempty"`
	Piridoxina *float64 `json:"piridoxina,omitempty"`
	// TBCA specific
	CarboidratoDisponivel *float64 `json:"carboidratoDisponivel,omitempty"`
	Alcool                *float64 `json:"alcool,omitempty"`
	AcidosGraxosTrans     *float64 `json:"acidosGraxosTrans,omitempty"`
	Selenio               *float64 `json:"selenio,omitempty"`
	VitaminaD             *float64 `json:"vitaminaD,omitempty"`
	AlfaTocoferol         *float64 `json:"alfaTocoferol,omitempty"`
	VitaminaB6            *float64 `json:"vitaminaB6,omitempty"`
	VitaminaB12           *float64 `json:"vitaminaB12,omitempty"`
	EquivalenteFolato     *float64 `json:"equivalenteFolato,omitempty"`
	SalAdicao             *float64 `json:"salAdicao,omitempty"`
	AcucarAdicao          *float64 `json:"acucarAdicao,omitempty"`
	NomeCientifico        string   `json:"nomeCientifico,omitempty"`
}

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func parseServingMeasures(data sql.NullString) []ServingMeasure {
	if !data.Valid || data.String == "" {
		return []ServingMeasure{}
	}

	var measures []ServingMeasure
	if err := json.Unmarshal([]byte(data.String), &measures); err != nil {
		return []ServingMeasure{}
	}
	return measures
}

func scanNullFloat(ns sql.NullFloat64) *float64 {
	if ns.Valid {
		return &ns.Float64
	}
	return nil
}

// parseFloatString handles strings like "NA", "-", "Tr" (trace) as null values
func parseFloatString(s sql.NullString) *float64 {
	if !s.Valid || s.String == "" || s.String == "NA" || s.String == "-" || s.String == "Tr" || s.String == "*" {
		return nil
	}
	// Replace comma with dot for Brazilian number format
	str := strings.Replace(s.String, ",", ".", -1)
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil
	}
	return &val
}

func (db *DB) GetTacoFoods() ([]Food, error) {
	query := `
		SELECT 
			Código,
			Nome,
			"Outras medidas",
			Grupo,
			"Energia |kJ|",
			"Energia |kcal|",
			"Umidade |g|",
			"Carboidrato total |g|",
			"Proteína |g|",
			"Lipídios |g|",
			"Fibra alimentar |g|",
			"Cinzas |g|",
			"Colesterol |mg|",
			"Ácidos graxos saturados |g|",
			"Ácidos graxos monoinsaturados |g|",
			"Ácidos graxos poliinsaturados |g|",
			"Cálcio |mg|",
			"Ferro |mg|",
			"Sódio |mg|",
			"Magnésio |mg|",
			"Fósforo |mg|",
			"Potássio |mg|",
			"Manganês |mg|",
			"Zinco |mg|",
			"Retinol |mcg|",
			"Cobre |mg|",
			"Vitamina A (RE) |mcg|",
			"Vitamina A (RAE) |mcg|",
			"Tiamina |mg|",
			"Riboflavina |mg|",
			"Piridoxina |mg|",
			"Niacina |mg|",
			"Vitamina C |mg|"
		FROM taco
		ORDER BY Nome
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query taco: %w", err)
	}
	defer rows.Close()

	var foods []Food
	for rows.Next() {
		var (
			codigo        int
			nome          string
			outrasMedidas sql.NullString
			grupo         sql.NullString
			energiaKJ     sql.NullString
			energiaKcal   sql.NullString
			umidade       sql.NullString
			carboidrato   sql.NullString
			proteina      sql.NullString
			lipidios      sql.NullString
			fibra         sql.NullString
			cinzas        sql.NullString
			colesterol    sql.NullString
			agSaturados   sql.NullString
			agMono        sql.NullString
			agPoli        sql.NullString
			calcio        sql.NullString
			ferro         sql.NullString
			sodio         sql.NullString
			magnesio      sql.NullString
			fosforo       sql.NullString
			potassio      sql.NullString
			manganes      sql.NullString
			zinco         sql.NullString
			retinol       sql.NullString
			cobre         sql.NullString
			vitARE        sql.NullString
			vitARAE       sql.NullString
			tiamina       sql.NullString
			riboflavina   sql.NullString
			piridoxina    sql.NullString
			niacina       sql.NullString
			vitC          sql.NullString
		)

		err := rows.Scan(
			&codigo, &nome, &outrasMedidas, &grupo,
			&energiaKJ, &energiaKcal, &umidade, &carboidrato,
			&proteina, &lipidios, &fibra, &cinzas,
			&colesterol, &agSaturados, &agMono, &agPoli,
			&calcio, &ferro, &sodio, &magnesio,
			&fosforo, &potassio, &manganes, &zinco,
			&retinol, &cobre, &vitARE, &vitARAE,
			&tiamina, &riboflavina, &piridoxina, &niacina, &vitC,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		grupoStr := ""
		if grupo.Valid {
			grupoStr = grupo.String
		}

		food := Food{
			Codigo:                      fmt.Sprintf("%d", codigo),
			Nome:                        nome,
			Slug:                        slug.Make(nome),
			OutrasMedidas:               parseServingMeasures(outrasMedidas),
			Grupo:                       grupoStr,
			Source:                      "taco",
			EnergiaKJ:                   parseFloatString(energiaKJ),
			EnergiaKcal:                 parseFloatString(energiaKcal),
			Umidade:                     parseFloatString(umidade),
			CarboidratoTotal:            parseFloatString(carboidrato),
			Proteina:                    parseFloatString(proteina),
			Lipidios:                    parseFloatString(lipidios),
			FibraAlimentar:              parseFloatString(fibra),
			Cinzas:                      parseFloatString(cinzas),
			Colesterol:                  parseFloatString(colesterol),
			AcidosGraxosSaturados:       parseFloatString(agSaturados),
			AcidosGraxosMonoinsaturados: parseFloatString(agMono),
			AcidosGraxosPoliinsaturados: parseFloatString(agPoli),
			Calcio:                      parseFloatString(calcio),
			Ferro:                       parseFloatString(ferro),
			Sodio:                       parseFloatString(sodio),
			Magnesio:                    parseFloatString(magnesio),
			Fosforo:                     parseFloatString(fosforo),
			Potassio:                    parseFloatString(potassio),
			Manganes:                    parseFloatString(manganes),
			Zinco:                       parseFloatString(zinco),
			Retinol:                     parseFloatString(retinol),
			Cobre:                       parseFloatString(cobre),
			VitaminaARE:                 parseFloatString(vitARE),
			VitaminaARAE:                parseFloatString(vitARAE),
			Tiamina:                     parseFloatString(tiamina),
			Riboflavina:                 parseFloatString(riboflavina),
			Piridoxina:                  parseFloatString(piridoxina),
			Niacina:                     parseFloatString(niacina),
			VitaminaC:                   parseFloatString(vitC),
		}
		foods = append(foods, food)
	}

	return foods, nil
}

func (db *DB) GetTbcaFoods() ([]Food, error) {
	query := `
		SELECT 
			Código,
			Nome,
			"Outras medidas",
			"Nome científico",
			Grupo,
			"Energia |kJ|",
			"Energia |kcal|",
			"Umidade |g|",
			"Carboidrato total |g|",
			"Carboidrato disponível |g|",
			"Proteína |g|",
			"Lipídios |g|",
			"Fibra alimentar |g|",
			"Álcool |g|",
			"Cinzas |g|",
			"Colesterol |mg|",
			"Ácidos graxos saturados |g|",
			"Ácidos graxos monoinsaturados |g|",
			"Ácidos graxos poliinsaturados |g|",
			"Ácidos graxos trans |g|",
			"Cálcio |mg|",
			"Ferro |mg|",
			"Sódio |mg|",
			"Magnésio |mg|",
			"Fósforo |mg|",
			"Potássio |mg|",
			"Manganês |mg|",
			"Zinco |mg|",
			"Cobre |mg|",
			"Selênio |mcg|",
			"Vitamina A (RE) |mcg|",
			"Vitamina A (RAE) |mcg|",
			"Vitamina D |mcg|",
			"Alfa-tocoferol (Vitamina E) |mg|",
			"Tiamina |mg|",
			"Riboflavina |mg|",
			"Niacina |mg|",
			"Vitamina B6 |mg|",
			"Vitamina B12 |mcg|",
			"Vitamina C |mg|",
			"Equivalente de folato |mcg|",
			"Sal de adição |g|",
			"Açúcar de adição |g|"
		FROM tbca
		ORDER BY Nome
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tbca: %w", err)
	}
	defer rows.Close()

	var foods []Food
	for rows.Next() {
		var (
			codigo         string
			nome           string
			outrasMedidas  sql.NullString
			nomeCientifico sql.NullString
			grupo          sql.NullString
			energiaKJ      sql.NullString
			energiaKcal    sql.NullString
			umidade        sql.NullString
			carboidrato    sql.NullString
			carbDisp       sql.NullString
			proteina       sql.NullString
			lipidios       sql.NullString
			fibra          sql.NullString
			alcool         sql.NullString
			cinzas         sql.NullString
			colesterol     sql.NullString
			agSaturados    sql.NullString
			agMono         sql.NullString
			agPoli         sql.NullString
			agTrans        sql.NullString
			calcio         sql.NullString
			ferro          sql.NullString
			sodio          sql.NullString
			magnesio       sql.NullString
			fosforo        sql.NullString
			potassio       sql.NullString
			manganes       sql.NullString
			zinco          sql.NullString
			cobre          sql.NullString
			selenio        sql.NullString
			vitARE         sql.NullString
			vitARAE        sql.NullString
			vitD           sql.NullString
			alfaTocoferol  sql.NullString
			tiamina        sql.NullString
			riboflavina    sql.NullString
			niacina        sql.NullString
			vitB6          sql.NullString
			vitB12         sql.NullString
			vitC           sql.NullString
			eqFolato       sql.NullString
			salAdicao      sql.NullString
			acucarAdicao   sql.NullString
		)

		err := rows.Scan(
			&codigo, &nome, &outrasMedidas, &nomeCientifico, &grupo,
			&energiaKJ, &energiaKcal, &umidade, &carboidrato, &carbDisp,
			&proteina, &lipidios, &fibra, &alcool, &cinzas,
			&colesterol, &agSaturados, &agMono, &agPoli, &agTrans,
			&calcio, &ferro, &sodio, &magnesio, &fosforo,
			&potassio, &manganes, &zinco, &cobre, &selenio,
			&vitARE, &vitARAE, &vitD, &alfaTocoferol, &tiamina,
			&riboflavina, &niacina, &vitB6, &vitB12, &vitC,
			&eqFolato, &salAdicao, &acucarAdicao,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		grupoStr := ""
		if grupo.Valid {
			grupoStr = grupo.String
		}

		nomeCientificoStr := ""
		if nomeCientifico.Valid {
			nomeCientificoStr = nomeCientifico.String
		}

		food := Food{
			Codigo:                      codigo,
			Nome:                        nome,
			Slug:                        slug.Make(nome) + "-tbca",
			OutrasMedidas:               parseServingMeasures(outrasMedidas),
			Grupo:                       grupoStr,
			Source:                      "tbca",
			NomeCientifico:              nomeCientificoStr,
			EnergiaKJ:                   parseFloatString(energiaKJ),
			EnergiaKcal:                 parseFloatString(energiaKcal),
			Umidade:                     parseFloatString(umidade),
			CarboidratoTotal:            parseFloatString(carboidrato),
			CarboidratoDisponivel:       parseFloatString(carbDisp),
			Proteina:                    parseFloatString(proteina),
			Lipidios:                    parseFloatString(lipidios),
			FibraAlimentar:              parseFloatString(fibra),
			Alcool:                      parseFloatString(alcool),
			Cinzas:                      parseFloatString(cinzas),
			Colesterol:                  parseFloatString(colesterol),
			AcidosGraxosSaturados:       parseFloatString(agSaturados),
			AcidosGraxosMonoinsaturados: parseFloatString(agMono),
			AcidosGraxosPoliinsaturados: parseFloatString(agPoli),
			AcidosGraxosTrans:           parseFloatString(agTrans),
			Calcio:                      parseFloatString(calcio),
			Ferro:                       parseFloatString(ferro),
			Sodio:                       parseFloatString(sodio),
			Magnesio:                    parseFloatString(magnesio),
			Fosforo:                     parseFloatString(fosforo),
			Potassio:                    parseFloatString(potassio),
			Manganes:                    parseFloatString(manganes),
			Zinco:                       parseFloatString(zinco),
			Cobre:                       parseFloatString(cobre),
			Selenio:                     parseFloatString(selenio),
			VitaminaARE:                 parseFloatString(vitARE),
			VitaminaARAE:                parseFloatString(vitARAE),
			VitaminaD:                   parseFloatString(vitD),
			AlfaTocoferol:               parseFloatString(alfaTocoferol),
			Tiamina:                     parseFloatString(tiamina),
			Riboflavina:                 parseFloatString(riboflavina),
			Niacina:                     parseFloatString(niacina),
			VitaminaB6:                  parseFloatString(vitB6),
			VitaminaB12:                 parseFloatString(vitB12),
			VitaminaC:                   parseFloatString(vitC),
			EquivalenteFolato:           parseFloatString(eqFolato),
			SalAdicao:                   parseFloatString(salAdicao),
			AcucarAdicao:                parseFloatString(acucarAdicao),
		}
		foods = append(foods, food)
	}

	return foods, nil
}

func (db *DB) GetAllFoods(source string) ([]Food, error) {
	switch source {
	case "taco":
		return db.GetTacoFoods()
	case "tbca":
		return db.GetTbcaFoods()
	case "both":
		tacoFoods, err := db.GetTacoFoods()
		if err != nil {
			return nil, err
		}
		tbcaFoods, err := db.GetTbcaFoods()
		if err != nil {
			return nil, err
		}
		return append(tacoFoods, tbcaFoods...), nil
	default:
		return db.GetTacoFoods()
	}
}
