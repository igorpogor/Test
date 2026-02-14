using System;
using System.Collections.Generic;

namespace CSharpLab
{
    class Program
    {
        static double CalculateF(double x, double y)
        {
            double x2 = x * x;
            double baseValue = (x2 + y) * (x2 + y);
            double xy = x * y;

            if (xy > 0)
            {
                return baseValue - Math.Sqrt(x2 * y);
            }
            else if (xy < 0)
            {
                return baseValue - Math.Sqrt(Math.Abs(x2 * y));
            }
            else
            {
                return baseValue + 1;
            }
        }

        static double ReadDouble(string prompt)
        {
            while (true)
            {
                Console.Write(prompt);
                string input = Console.ReadLine();
                if (input == null)
                {
                    Environment.Exit(0);
                }
                input = input.Replace('.', ',');
                double result;
                if (double.TryParse(input, out result))
                {
                    return result;
                }
                Console.WriteLine("Ошибка ввода. Попробуйте снова.");
            }
        }

        static void GetUserInput(out double x, out double y, out double n)
        {
            x = ReadDouble("Введите x: ");
            y = ReadDouble("Введите y: ");
            n = 0;
            while (n <= 0)
            {
                n = ReadDouble("Введите шаг n (> 0): ");
                if (n <= 0)
                    Console.WriteLine("Шаг должен быть положительным числом.");
            }
        }

        static void CalculateResultsRecursive(double startX, double y, double n, double endX,
            int i, List<double> xValues, List<double> fValues)
        {
            double currentX = startX + i * n;

            if (currentX > endX + 1e-9)
                return;

            double f = CalculateF(currentX, y);
            xValues.Add(currentX);
            fValues.Add(f);

            CalculateResultsRecursive(startX, y, n, endX, i + 1, xValues, fValues);
        }

        static void CalculateResults(double x, double y, double n,
            out List<double> xValues, out List<double> fValues)
        {
            double endX = x + 10.0;
            xValues = new List<double>();
            fValues = new List<double>();

            CalculateResultsRecursive(x, y, n, endX, 0, xValues, fValues);
        }

        static void DisplayResults(List<double> xValues, List<double> fValues)
        {
            Console.WriteLine();
            Console.WriteLine("{0,-5} {1,-20} {2,-25}", "i", "x", "f(x, y)");
            Console.WriteLine(new string('-', 50));

            for (int i = 0; i < xValues.Count; i++)
            {
                Console.WriteLine("{0,-5} {1,-20:F6} {2,-25:F6}", i + 1, xValues[i], fValues[i]);
            }

            Console.WriteLine();
            Console.WriteLine("Результаты сохранены в массив из {0} элементов.", fValues.Count);
        }

        static void Main(string[] args)
        {
            double x, y, n;
            GetUserInput(out x, out y, out n);

            List<double> xValues, fValues;
            CalculateResults(x, y, n, out xValues, out fValues);

            DisplayResults(xValues, fValues);

            Console.WriteLine();
            Console.WriteLine("Нажмите любую клавишу для выхода...");
            Console.ReadKey();
        }
    }
}