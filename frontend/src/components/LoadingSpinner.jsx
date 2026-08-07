export default function LoadingSpinner({ className = '' }) {
  return (
    <div className={`flex items-center justify-center ${className}`.trim()} aria-label="Loading">
      <div className="h-8 w-8 animate-spin rounded-full border-4 border-[#948979]/30 border-t-[#948979]" />
    </div>
  );
}
