package sim

import (
	"math"
	"sort"
	"time"
)

// PerformanceStats are classic backtest metrics.
type PerformanceStats struct {
	TotalReturnPct    float64 `json:"total_return_pct"`
	CAGRPct           float64 `json:"cagr_pct"`
	SharpeRatio       float64 `json:"sharpe_ratio"`
	MaxDrawdownPct    float64 `json:"max_drawdown_pct"`
	VolatilityPct     float64 `json:"volatility_pct"`
	WinRatePct        float64 `json:"win_rate_pct"`
	ProfitFactor      float64 `json:"profit_factor"`
	NumTrades         int     `json:"num_trades"`
	AvgTradeReturnPct float64 `json:"avg_trade_return_pct"`
}

// ComputePerformance calculates statistics from equity curve and candles.
func ComputePerformance(equity []EquityPoint, candles []Candle) PerformanceStats {
	if len(equity) == 0 || len(candles) == 0 {
		return PerformanceStats{}
	}
	start := equity[0].Value
	end := equity[len(equity)-1].Value
	if start == 0 {
		start = 1
	}
	returns := make([]float64, 0, len(equity)-1)
	peak := equity[0].Value
	maxDD := 0.0
	for i := 1; i < len(equity); i++ {
		r := (equity[i].Value - equity[i-1].Value) / equity[i-1].Value
		returns = append(returns, r)
		if equity[i].Value > peak {
			peak = equity[i].Value
		}
		dd := (peak - equity[i].Value) / peak
		if dd > maxDD {
			maxDD = dd
		}
	}
	sort.Float64s(returns)
	var sum, sqSum float64
	for _, r := range returns {
		sum += r
		sqSum += r * r
	}
	n := float64(len(returns))
	avg := 0.0
	std := 0.0
	if n > 0 {
		avg = sum / n
		std = math.Sqrt(sqSum/n - avg*avg)
	}
	sharpe := 0.0
	if std != 0 {
		sharpe = (avg * math.Sqrt(252)) / std
		if math.IsNaN(sharpe) || math.IsInf(sharpe, 0) {
			sharpe = 0
		}
	}
	totalReturn := (end - start) / start
	days := float64(len(candles))
	if days == 0 {
		days = 1
	}
	cagr := math.Pow(1+totalReturn, 252.0/days) - 1
	if math.IsNaN(cagr) || math.IsInf(cagr, 0) {
		cagr = 0
	}
	return PerformanceStats{
		TotalReturnPct: totalReturn * 100,
		CAGRPct:        cagr * 100,
		SharpeRatio:    sharpe,
		MaxDrawdownPct: maxDD * 100,
		VolatilityPct:  std * 100 * math.Sqrt(252),
	}
}

// AnnualizedDays returns the number of trading days between two times.
func AnnualizedDays(start, end time.Time) float64 {
	d := end.Sub(start).Hours() / 24
	if d < 1 {
		d = 1
	}
	return d
}
