import { File, Globe } from "lucide-react";
import { formatDate } from "@/components/paged-object-table";
import { cn, humanReadableSize, truncateUrl } from "@/lib/utils";

interface BaseItem {
  id: string;
  date_uploaded: string;
}

interface DocumentItem extends BaseItem {
  type: 'document';
  name: string;
  document_size: number;
}

interface WebpageItem extends BaseItem {
  type: 'webpage';
  url: string;
}

type RecentItem = DocumentItem | WebpageItem;

interface RecentItemsTableProps {
  items: RecentItem[];
  emptyMessage: string;
  title: string;
}

export function getRecentItems<T extends BaseItem>(
  items: T[] | undefined,
  daysToLookBack: number = 7,
  maxItems: number = 10
): T[] {
  if (!items) return [];
  
  const cutoffDate = new Date();
  cutoffDate.setDate(cutoffDate.getDate() - daysToLookBack);
  
  return items
    .filter(item => new Date(item.date_uploaded) > cutoffDate)
    .sort((a, b) => new Date(b.date_uploaded).getTime() - new Date(a.date_uploaded).getTime())
    .slice(0, maxItems);
}

export function RecentItemsTable({ items, emptyMessage, title }: RecentItemsTableProps) {
  return (
    <div className="border rounded-xl bg-card">
      <div className="overflow-hidden rounded-xl">
        <div className="border-b bg-table-header">
          <h3 className="p-4 font-semibold text-primary">
            {title}
          </h3>
        </div>

        <div className="divide-y divide-border">
          {items.length === 0 ? (
            <div className="p-4 text-center bg-table-row-odd text-foreground">
              {emptyMessage}
            </div>
          ) : (
            <div className="divide-y divide-border">
              {items.map((item, index) => (
                <div
                  key={item.id}
                  className={cn(
                    "p-4 flex items-center transition-colors",
                    index % 2 === 0 
                      ? "bg-table-row-odd" 
                      : "bg-table-row-even",
                    "hover:bg-muted/50"
                  )}
                >
                  <div className="flex items-center gap-3">
                    {item.type === 'document' ? (
                      <File className="size-5 text-muted-foreground" />
                    ) : (
                      <Globe className="size-5 text-muted-foreground" />
                    )}
                    <div>
                      {item.type === 'document' ? (
                        <>
                          <p className="font-medium text-foreground">{item.name}</p>
                          <p className="text-sm text-muted-foreground">
                            {humanReadableSize(item.document_size)} • {formatDate(item.date_uploaded)}
                          </p>
                        </>
                      ) : (
                        <>
                          <p className="font-medium text-foreground font-mono">{truncateUrl(item.url)}</p>
                          <p className="text-sm text-muted-foreground">
                            {formatDate(item.date_uploaded)}
                          </p>
                        </>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
} 