import type { Food } from "../types/food";

// Use port 8085 for the API (Go server must be running during build)
// Can be overridden via API_URL environment variable
const API_URL =
  (typeof process !== "undefined" && process.env?.API_URL) ||
  "http://localhost:8085";

export async function fetchFoods(
  source: "taco" | "tbca" | "both" = "both"
): Promise<Food[]> {
  const response = await fetch(`${API_URL}/api/foods?source=${source}`);
  if (!response.ok) {
    throw new Error("Failed to fetch foods");
  }
  return response.json();
}

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
