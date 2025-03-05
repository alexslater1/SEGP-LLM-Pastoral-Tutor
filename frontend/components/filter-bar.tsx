import { cn } from "@/lib/utils";
import { Button } from "./ui/button";
import { LucideIcon } from "lucide-react";

type FilterOption = {
  id: string;
  label: string;
  icon: LucideIcon;
  theme: 'default' | 'success' | 'destructive';
}

export function FilterBar<T extends string>({
  filter,
  setFilter,
  options
}: {
  filter: T;
  setFilter: (filter: T) => void;
  options: FilterOption[];
}) {
  const getThemeStyles = (theme: FilterOption['theme']) => {
    switch (theme) {
      case 'success':
        return "bg-success/10 border-success text-success hover:bg-success/20";
      case 'destructive':
        return "bg-destructive/10 border-destructive text-destructive hover:bg-destructive/20";
      default:
        return "bg-primary/10 border-primary text-primary hover:bg-primary/20";
    }
  };

  return (
    <div className="absolute right-6 top-20">
      <div className="relative flex items-center gap-2 bg-background rounded-lg p-1 shadow-sm">
        {options.map((option) => (
          <Button
            key={option.id}
            variant="outline"
            className={cn(
              "flex items-center gap-2",
              filter === option.id && getThemeStyles(option.theme)
            )}
            onClick={() => setFilter(option.id as T)}
          >
            <option.icon size={16} />
            {option.label}
          </Button>
        ))}
      </div>
    </div>
  );
} 