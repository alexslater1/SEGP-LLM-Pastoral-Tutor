import { cn } from "@/lib/utils";
import { Button } from "./ui/button";

export function PageSelectionBar({pageNumber, setPage, pageNames}: {pageNumber: number, setPage: (page: number) => void, pageNames: string[]}) {
  return (
    <div className="absolute right-6 top-16">
      <div className="relative flex items-center bg-card rounded-lg p-1 border shadow-sm">
        <ButtonBackground pageNumber={pageNumber} pages={pageNames} animate={true} />

        {pageNames.map((page, index) => (
          <PageSelectionBarButton
            key={"selectionButton" + index}
            pageName={page}
            pageNumbers={[pageNumber]}
            setPage={setPage}
            unsetPage={() => {}}
            index={index}
          />
        ))}
      </div>
    </div>
  )
}

export function PageSelectionBarMultiple({pageNumbers, setPages, pageNames, allowNoneSelected}: 
  {pageNumbers: number[], setPages: (page: number[]) => void, pageNames: string[], allowNoneSelected: boolean}) {
  
  const appendPageNumber = (pageNumber: number) => {
    setPages([...pageNumbers, pageNumber]);
  }

  const removePageNumber = (pageNumber: number) => {
    if (allowNoneSelected || pageNumbers.length > 1) {
      setPages(pageNumbers.filter((number) => number !== pageNumber));
    }
  }

  return (
    <div className="absolute right-6 top-16">
      <div className="relative flex items-center bg-card rounded-lg p-1 border shadow-sm">
        {pageNumbers.map((pageNumber, index) => (
          <ButtonBackground 
            key={"buttonBackground" + index}
            pageNumber={pageNumber} 
            pages={pageNames} 
            animate={false} 
          />
        ))}

        {pageNames.map((page, index) => (
          <PageSelectionBarButton
            key={"selectionButtonMultiple" + index}
            pageName={page}
            pageNumbers={pageNumbers}
            setPage={appendPageNumber}
            unsetPage={removePageNumber}
            index={index}
          />
        ))}
      </div>
    </div>
  )
}

function ButtonBackground({pageNumber, pages, animate}: {pageNumber: number, pages: string[], animate: boolean}) {
  return (
    <div 
      className={cn(
        "absolute h-[85%] top-[7.5%] bg-primary/10 rounded-md ",
        animate && "transition-all duration-300 ease-out"
      )}
      style={{
        left: 'calc(' + (100 / pages.length) * pageNumber + '% + 3px)',
        width: 'calc(' + (100 / pages.length) + '% - 6px)',
      }}
    />
  )
}

function PageSelectionBarButton({pageName, pageNumbers, setPage, unsetPage, index}: 
  {pageName: string, pageNumbers: number[], setPage: (page: number) => void, unsetPage: (page: number) => void, index: number}) {
  return (
    <Button
      onClick={pageNumbers.includes(index) ? () => unsetPage(index) : () => setPage(index)}
      variant="ghost"
      className={cn(
        "relative px-2 z-10 transition-colors duration-300 w-24",
        pageNumbers.includes(index) ? "font-bold text-primary hover:bg-transparent" : "text-muted-foreground"
      )}
    >
      {pageName}
    </Button>
  )
}


