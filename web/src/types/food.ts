export interface ServingMeasure {
  label: string;
  value: string;
}

export interface Food {
  codigo: string;
  nome: string;
  slug: string;
  outrasMedidas: ServingMeasure[];
  grupo: string;
  source: "taco" | "tbca";
  energiaKJ: number | null;
  energiaKcal: number | null;
  umidade: number | null;
  carboidratoTotal: number | null;
  proteina: number | null;
  lipidios: number | null;
  fibraAlimentar: number | null;
  cinzas: number | null;
  colesterol: number | null;
  acidosGraxosSaturados: number | null;
  acidosGraxosMonoinsaturados: number | null;
  acidosGraxosPoliinsaturados: number | null;
  calcio: number | null;
  ferro: number | null;
  sodio: number | null;
  magnesio: number | null;
  fosforo: number | null;
  potassio: number | null;
  manganes: number | null;
  zinco: number | null;
  cobre: number | null;
  vitaminaARE: number | null;
  vitaminaARAE: number | null;
  tiamina: number | null;
  riboflavina: number | null;
  niacina: number | null;
  vitaminaC: number | null;
  // TACO specific
  retinol?: number | null;
  piridoxina?: number | null;
  // TBCA specific
  carboidratoDisponivel?: number | null;
  alcool?: number | null;
  acidosGraxosTrans?: number | null;
  selenio?: number | null;
  vitaminaD?: number | null;
  alfaTocoferol?: number | null;
  vitaminaB6?: number | null;
  vitaminaB12?: number | null;
  equivalenteFolato?: number | null;
  salAdicao?: number | null;
  acucarAdicao?: number | null;
  nomeCientifico?: string;
}

export interface NutrientInfo {
  key: keyof Food;
  label: string;
  unit: string;
  category:
    | "principal"
    | "fibras"
    | "lipidios"
    | "minerais"
    | "vitaminas"
    | "outros";
}

export const nutrientsList: NutrientInfo[] = [
  // Principal
  { key: "energiaKcal", label: "Energia", unit: "kcal", category: "principal" },
  {
    key: "carboidratoTotal",
    label: "Carboidratos",
    unit: "g",
    category: "principal",
  },
  { key: "proteina", label: "Proteínas", unit: "g", category: "principal" },
  {
    key: "lipidios",
    label: "Gorduras Totais",
    unit: "g",
    category: "principal",
  },
  {
    key: "fibraAlimentar",
    label: "Fibra Alimentar",
    unit: "g",
    category: "principal",
  },

  // Lipídios
  {
    key: "acidosGraxosSaturados",
    label: "Gorduras Saturadas",
    unit: "g",
    category: "lipidios",
  },
  {
    key: "acidosGraxosMonoinsaturados",
    label: "Gorduras Monoinsaturadas",
    unit: "g",
    category: "lipidios",
  },
  {
    key: "acidosGraxosPoliinsaturados",
    label: "Gorduras Poliinsaturadas",
    unit: "g",
    category: "lipidios",
  },
  {
    key: "acidosGraxosTrans",
    label: "Gorduras Trans",
    unit: "g",
    category: "lipidios",
  },

  // Minerais
  { key: "sodio", label: "Sódio", unit: "mg", category: "minerais" },
  { key: "calcio", label: "Cálcio", unit: "mg", category: "minerais" },
  { key: "ferro", label: "Ferro", unit: "mg", category: "minerais" },
  { key: "potassio", label: "Potássio", unit: "mg", category: "minerais" },
  { key: "magnesio", label: "Magnésio", unit: "mg", category: "minerais" },
  { key: "fosforo", label: "Fósforo", unit: "mg", category: "minerais" },
  { key: "zinco", label: "Zinco", unit: "mg", category: "minerais" },
  { key: "manganes", label: "Manganês", unit: "mg", category: "minerais" },
  { key: "cobre", label: "Cobre", unit: "mg", category: "minerais" },
  { key: "selenio", label: "Selênio", unit: "mcg", category: "minerais" },

  // Vitaminas
  { key: "vitaminaC", label: "Vitamina C", unit: "mg", category: "vitaminas" },
  {
    key: "vitaminaARAE",
    label: "Vitamina A (RAE)",
    unit: "mcg",
    category: "vitaminas",
  },
  {
    key: "vitaminaARE",
    label: "Vitamina A (RE)",
    unit: "mcg",
    category: "vitaminas",
  },
  { key: "vitaminaD", label: "Vitamina D", unit: "mcg", category: "vitaminas" },
  {
    key: "alfaTocoferol",
    label: "Vitamina E",
    unit: "mg",
    category: "vitaminas",
  },
  { key: "tiamina", label: "Tiamina (B1)", unit: "mg", category: "vitaminas" },
  {
    key: "riboflavina",
    label: "Riboflavina (B2)",
    unit: "mg",
    category: "vitaminas",
  },
  { key: "niacina", label: "Niacina (B3)", unit: "mg", category: "vitaminas" },
  {
    key: "piridoxina",
    label: "Piridoxina (B6)",
    unit: "mg",
    category: "vitaminas",
  },
  {
    key: "vitaminaB6",
    label: "Vitamina B6",
    unit: "mg",
    category: "vitaminas",
  },
  {
    key: "vitaminaB12",
    label: "Vitamina B12",
    unit: "mcg",
    category: "vitaminas",
  },
  {
    key: "equivalenteFolato",
    label: "Folato",
    unit: "mcg",
    category: "vitaminas",
  },
  { key: "retinol", label: "Retinol", unit: "mcg", category: "vitaminas" },

  // Outros
  { key: "colesterol", label: "Colesterol", unit: "mg", category: "outros" },
  { key: "alcool", label: "Álcool", unit: "g", category: "outros" },
  {
    key: "carboidratoDisponivel",
    label: "Carboidrato Disponível",
    unit: "g",
    category: "outros",
  },
  { key: "salAdicao", label: "Sal de Adição", unit: "g", category: "outros" },
  {
    key: "acucarAdicao",
    label: "Açúcar de Adição",
    unit: "g",
    category: "outros",
  },
];

export function formatNutrientValue(
  value: number | null | undefined,
  multiplier: number = 1
): string {
  if (value === null || value === undefined) {
    return "-";
  }
  const calculated = value * multiplier;
  if (calculated === 0) return "0";
  if (calculated < 0.01) return "<0,01";
  if (calculated < 1) return calculated.toFixed(2).replace(".", ",");
  if (calculated < 10) return calculated.toFixed(1).replace(".", ",");
  return Math.round(calculated).toString();
}
